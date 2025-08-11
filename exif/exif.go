// Package exif implements decoding of EXIF data as defined in the EXIF 2.2
// specification (http://www.exif.org/Exif2-2.PDF).
package exif

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	heif "github.com/jdeng/goheif"
	"github.com/rpajarola/exiftools/models"
	"github.com/rpajarola/exiftools/tiff"
)

// DecodeOptions provides configuration options for EXIF decoding
type DecodeOptions struct {
	KeepUnknownTags bool // Keep unknown tags (default: false)
	MaxExifSize     int  // maximum size of exif data (default: 4MB)
}

const (
	// JPEG APP1 section marker
	jpegAPP1 = 0xE1

	// limits the exif size to avoid trying to read corrupted lengths
	// and parsing potentially gigabytes of exif
	exifLengthCutoff = 4 * 1024 * 1024
)

// Exif provides access to decoded EXIF metadata fields and values.
type Exif struct {
	Tiff *tiff.Tiff
	Fields map[models.FieldName]*tiff.Tag
	Raw  []byte
	opts DecodeOptions
}

func New(tif *tiff.Tiff, raw []byte, opts *DecodeOptions) *Exif {
	if opts == nil {
		opts = &DecodeOptions{}
	}
	return &Exif{
		Fields: map[models.FieldName]*tiff.Tag{},
		Tiff: tif,
		Raw:  raw,
		opts: *opts,
	}
}

// Decode parses EXIF data from r (a TIFF, JPEG, or raw EXIF block)
// and returns a queryable Exif object. After the EXIF data section is
// called and the TIFF structure is decoded, each registered parser is
// called (in order of registration). If one parser returns an error,
// decoding terminates and the remaining parsers are not called.
//
// The error can be inspected with functions such as IsCriticalError
// to determine whether the returned object might still be usable.
func Decode(r io.Reader) (*Exif, error) {
	return DecodeWithOptions(r, &DecodeOptions{KeepUnknownTags: false})
}

// fileType represents the detected file format
type fileType int

const (
	fileTypeTIFF fileType = iota
	fileTypeRawExif
	fileTypeHEIF
	fileTypeJPEG
)

// detectFileType examines the header to determine the file format
func detectFileType(header []byte) fileType {
	switch string(header[0:4]) {
	case "II*\x00", "MM\x00*":
		// TIFF - Little/Big endian
		return fileTypeTIFF
	case "Exif":
		return fileTypeRawExif
	default:
		if string(header[4:]) == "ftyp" {
			return fileTypeHEIF
		}
		// Assume JPEG
		return fileTypeJPEG
	}
}

// processHEIFFile extracts EXIF data from HEIF/HEIC files
func processHEIFFile(r io.Reader) (io.Reader, error) {
	// For HEIF files, we need a ReaderAt interface
	// Try to use a more memory-efficient approach when possible
	if ra, ok := r.(io.ReaderAt); ok {
		// If we already have a ReaderAt, use it directly
		xf, err := heif.ExtractExif(ra)
		if err != nil {
			return nil, fmt.Errorf("exif: unable to extract exif from heif/heic file: %w", err)
		}
		return bytes.NewReader(xf), nil
	} else {
		// Fallback: read into memory (necessary for ReaderAt interface)
		data, err := io.ReadAll(r)
		if err != nil {
			return nil, fmt.Errorf("exif: unable to read heif/heic file: %w", err)
		}
		ra := bytes.NewReader(data)
		xf, err := heif.ExtractExif(ra)
		if err != nil {
			return nil, fmt.Errorf("exif: unable to extract exif from heif/heic file: %w", err)
		}
		return bytes.NewReader(xf), nil
	}
}

// processRawExifFile validates and processes raw EXIF data
func processRawExifFile(r io.Reader) error {
	var header [6]byte
	if _, err := io.ReadFull(r, header[:]); err != nil {
		return fmt.Errorf("exif: unexpected raw exif header read error")
	}
	if got, want := string(header[:]), "Exif\x00\x00"; got != want {
		return fmt.Errorf("exif: unexpected raw exif header; got %q, want %q", got, want)
	}
	return nil
}

// processTIFFFile decodes TIFF data and returns the bytes reader and tiff structure
func processTIFFFile(r io.Reader) (*bytes.Reader, *tiff.Tiff, error) {
	// Functions below need the IFDs from the TIFF data to be stored in a
	// *bytes.Reader.  We use TeeReader to get a copy of the bytes as a
	// side-effect of tiff.Decode() doing its work.
	b := &bytes.Buffer{}
	tr := io.TeeReader(r, b)
	tif, err := tiff.Decode(tr)
	if err != nil {
		return nil, nil, err
	}
	return bytes.NewReader(b.Bytes()), tif, nil
}

// processJPEGFile extracts EXIF data from JPEG APP1 section
func processJPEGFile(r io.Reader) (*bytes.Reader, *tiff.Tiff, error) {
	// Locate the JPEG APP1 header.
	sec, err := newAppSec(jpegAPP1, r)
	if err != nil {
		return nil, nil, err
	}
	// Strip away EXIF header.
	er, err := sec.exifReader()
	if err != nil {
		return nil, nil, err
	}
	tif, err := tiff.Decode(er)
	if err != nil {
		return nil, nil, err
	}
	return er, tif, nil
}

// DecodeWithOptions parses EXIF data with the provided options.
func DecodeWithOptions(r io.Reader, opts *DecodeOptions) (*Exif, error) {
	if opts == nil {
		opts = &DecodeOptions{}
	}
	if opts.MaxExifSize <= 0 {
		opts.MaxExifSize = exifLengthCutoff
	}

	// Read header to detect file type
	header := make([]byte, 8)
	n, err := io.ReadFull(r, header)
	if err != nil {
		return nil, fmt.Errorf("exif: error reading 8 byte header, got %d, %v", n, err)
	}

	// Detect the file type
	fType := detectFileType(header)

	// Put the header bytes back into the reader.
	r = io.MultiReader(bytes.NewReader(header), r)

	var (
		er  *bytes.Reader
		tif *tiff.Tiff
	)

	// Process based on file type
	switch fType {
	case fileTypeHEIF:
		r, err = processHEIFFile(r)
		if err != nil {
			return nil, err
		}
		fallthrough
	case fileTypeRawExif:
		if fType == fileTypeRawExif {
			if err = processRawExifFile(r); err != nil {
				return nil, err
			}
		}
		fallthrough
	case fileTypeTIFF:
		er, tif, err = processTIFFFile(r)
	case fileTypeJPEG:
		er, tif, err = processJPEGFile(r)
	}

	if err != nil {
		return nil, decodeError{cause: err}
	}

	// Read raw EXIF data
	er.Seek(0, 0)
	raw, err := io.ReadAll(er)
	if err != nil {
		return nil, decodeError{cause: err}
	}

	// Build EXIF structure
	x := New(tif, raw, opts)

	// Run parsers
	for i, p := range parsers {
		if err := p.Parse(x); err != nil {
			if _, ok := err.(tiffErrors); ok {
				return x, err
			}
			// This should never happen, as Parse always returns a tiffError
			// for now, but that could change.
			return x, fmt.Errorf("exif: parser %v failed (%v)", i, err)
		}
	}

	return x, nil
}

// LoadTags loads tags into the available fields from the tiff Directory
// using the given tagid-fieldname mapping.  Used to load makernote and
// other meta-data.  If showMissing is true, tags in d that are not in the
// fieldMap will be loaded with the FieldName UnknownPrefix followed by the
// tag ID (in hex format).
func (x *Exif) LoadTags(d *tiff.Dir, fieldMap map[uint16]models.FieldName, showMissing bool) {
	for _, tag := range d.Tags {
		name := fieldMap[tag.Id]
		if name == "" {
			if !showMissing {
				continue
			}
			name = models.FieldName(fmt.Sprintf("%v%x", models.UnknownPrefix, tag.Id))
		}
		x.Fields[name] = tag
	}
}

// Get retrieves the EXIF tag for the given field name.
//
// If the tag is not known or not present, an error is returned. If the
// tag name is known, the error will be a TagNotPresentError.
func (x *Exif) Get(name models.FieldName) (*tiff.Tag, error) {
	if tg, ok := x.Fields[name]; ok {
		return tg, nil
	}
	return nil, TagNotPresentError(name)
}

// Update -- Illegal
func (x *Exif) Update(name models.FieldName, tag *tiff.Tag) error {
	if _, ok := x.Fields[name]; ok {
		x.Fields[name] = tag
		return nil
	}
	return TagNotPresentError(name)
}

// Walker is the interface used to traverse all fields of an Exif object.
type Walker interface {
	// Walk is called for each non-nil EXIF field. Returning a non-nil
	// error aborts the walk/traversal.
	Walk(name models.FieldName, tag *tiff.Tag) error
}

// Walk calls the Walk method of w with the name and tag for every non-nil
// EXIF field.  If w aborts the walk with an error, that error is returned.
func (x *Exif) Walk(w Walker) error {
	for name, tag := range x.Fields {
		if err := w.Walk(name, tag); err != nil {
			return err
		}
	}
	return nil
}

// DateTime returns the EXIF's "DateTimeOriginal" field, which
// is the creation time of the photo. If not found, it tries
// the "DateTime" (which is meant as the modtime) instead.
// The error will be TagNotPresentErr if none of those tags
// were found, or a generic error if the tag value was
// not a string, or the error returned by time.Parse.
//
// If the EXIF lacks timezone information or GPS time, the returned
// time's Location will be time.Local.
func (x *Exif) DateTime(fields ...models.FieldName) (time.Time, error) {
	var dt time.Time
	var tag *tiff.Tag
	var err error
	if len(fields) == 0 {
		fields = append(fields, models.DateTimeOriginal, models.DateTime)
	}
	for _, f := range fields {
		tag, err = x.Get(f)
		if err == nil {
			break
		}
	}
	if tag == nil {
		return dt, errors.New("DateTime not found")
	}
	if tag.Format() != tiff.StringVal {
		return dt, errors.New("DateTime[Original] not in string format")
	}

	exifTimeLayout := "2006:01:02 15:04:05"
	dateStr := strings.TrimRight(string(tag.Val), "\x00")
	subSecTag, err := x.Get(models.SubSecTimeOriginal)
	if err == nil {
		subSec, err := subSecTag.StringVal()
		if err == nil {
			if len(subSec) == 2 {
				exifTimeLayout = "2006:01:02 15:04:05.99"
			} else if len(subSec) == 3 {
				exifTimeLayout = "2006:01:02 15:04:05.999"
			}
			dateStr = fmt.Sprintf("%v.%v", dateStr, subSec)
		}
	}

	// TODO(bradfitz,mpl): look for timezone offset, GPS time, etc.
	//timeZone := time.Local
	//if tz, _ := x.TimeZone(); tz != nil {
	//	timeZone = tz
	//}
	return time.Parse(exifTimeLayout, dateStr)
	//return time.ParseInLocation(exifTimeLayout, dateStr, timeZone)
}

// TimeZone -
func (x *Exif) TimeZone() (*time.Location, error) {
	// TODO: parse more timezone fields (e.g. Nikon WorldTime).
	timeInfo, err := x.Get("Canon.TimeInfo")
	if err != nil {
		return nil, err
	}
	if timeInfo.Count < 2 {
		return nil, errors.New("Canon.TimeInfo does not contain timezone")
	}
	offsetMinutes, err := timeInfo.Int(1)
	if err != nil {
		return nil, err
	}
	return time.FixedZone("", offsetMinutes*60), nil
}

func ratFloat(num, dem int64) float64 {
	return float64(num) / float64(dem)
}

// Tries to parse a Geo degrees value from a string as it was found in some
// EXIF data.
// Supported formats so far:
//   - "52,00000,50,00000,34,01180" ==> 52 deg 50'34.0118"
//     Probably due to locale the comma is used as decimal mark as well as the
//     separator of three floats (degrees, minutes, seconds)
//     http://en.wikipedia.org/wiki/Decimal_mark#Hindu.E2.80.93Arabic_numeral_system
//   - "52.0,50.0,34.01180" ==> 52deg50'34.0118"
//   - "52,50,34.01180"     ==> 52deg50'34.0118"
func parseTagDegreesString(s string) (float64, error) {
	const unparsableErrorFmt = "unknown coordinate format: %s"
	isSplitRune := func(c rune) bool {
		return c == ',' || c == ';'
	}
	parts := strings.FieldsFunc(s, isSplitRune)
	var degrees, minutes, seconds float64
	var err error
	switch len(parts) {
	case 6:
		degrees, err = strconv.ParseFloat(parts[0]+"."+parts[1], 64)
		if err != nil {
			return 0.0, fmt.Errorf(unparsableErrorFmt, s)
		}
		minutes, err = strconv.ParseFloat(parts[2]+"."+parts[3], 64)
		if err != nil {
			return 0.0, fmt.Errorf(unparsableErrorFmt, s)
		}
		minutes = math.Copysign(minutes, degrees)
		seconds, err = strconv.ParseFloat(parts[4]+"."+parts[5], 64)
		if err != nil {
			return 0.0, fmt.Errorf(unparsableErrorFmt, s)
		}
		seconds = math.Copysign(seconds, degrees)
	case 3:
		degrees, err = strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return 0.0, fmt.Errorf(unparsableErrorFmt, s)
		}
		minutes, err = strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0.0, fmt.Errorf(unparsableErrorFmt, s)
		}
		minutes = math.Copysign(minutes, degrees)
		seconds, err = strconv.ParseFloat(parts[2], 64)
		if err != nil {
			return 0.0, fmt.Errorf(unparsableErrorFmt, s)
		}
		seconds = math.Copysign(seconds, degrees)
	default:
		return 0.0, fmt.Errorf(unparsableErrorFmt, s)
	}
	return degrees + minutes/60.0 + seconds/3600.0, nil
}

func parse3Rat2(tag *tiff.Tag) ([3]float64, error) {
	v := [3]float64{}
	for i := range v {
		num, den, err := tag.Rat2(i)
		if err != nil {
			return v, err
		}
		v[i] = ratFloat(num, den)
		if tag.Count < uint32(i+2) {
			break
		}
	}
	return v, nil
}

func tagDegrees(tag *tiff.Tag) (float64, error) {
	switch tag.Format() {
	case tiff.RatVal:
		// The usual case, according to the Exif spec
		// (http://www.kodak.com/global/plugins/acrobat/en/service/digCam/exifStandard2.pdf,
		// sec 4.6.6, p. 52 et seq.)
		v, err := parse3Rat2(tag)
		if err != nil {
			return 0.0, err
		}
		return v[0] + v[1]/60 + v[2]/3600.0, nil
	case tiff.StringVal:
		// Encountered this weird case with a panorama picture taken with a HTC phone
		s, err := tag.StringVal()
		if err != nil {
			return 0.0, err
		}
		return parseTagDegreesString(s)
	default:
		// don't know how to parse value, give up
		return 0.0, fmt.Errorf("malformed EXIF Tag Degrees")
	}
}

// parseGPSCoordinate extracts and parses a GPS coordinate (latitude or longitude)
// from the given coordinate and reference tags.
func (x *Exif) parseGPSCoordinate(coordTag, refTag models.FieldName, coordName string) (float64, error) {
	// Get the coordinate tag
	tag, err := x.Get(coordTag)
	if err != nil {
		return 0, err
	}

	// Get the reference tag (N/S for latitude, E/W for longitude)
	refTagVal, err := x.Get(refTag)
	if err != nil {
		return 0, err
	}

	// Parse the degrees value
	coord, err := tagDegrees(tag)
	if err != nil {
		return 0, fmt.Errorf("cannot parse %s: %v", coordName, err)
	}

	// Apply the reference direction
	ref, err := refTagVal.StringVal()
	if err != nil {
		return 0, fmt.Errorf("cannot parse %s reference: %v", coordName, err)
	}

	// Apply negative sign for South/West
	if ref == "S" || ref == "W" {
		coord *= -1.0
	}

	return coord, nil
}

// parseGPSLatitude extracts and parses GPS latitude from EXIF data.
func (x *Exif) parseGPSLatitude() (float64, error) {
	return x.parseGPSCoordinate(models.FieldName("GPSLatitude"), models.FieldName("GPSLatitudeRef"), "latitude")
}

// parseGPSLongitude extracts and parses GPS longitude from EXIF data.
func (x *Exif) parseGPSLongitude() (float64, error) {
	return x.parseGPSCoordinate(models.FieldName("GPSLongitude"), models.FieldName("GPSLongitudeRef"), "longitude")
}

// LatLong returns the latitude and longitude of the photo and
// whether it was present.
func (x *Exif) LatLong() (lat, long float64, err error) {
	lat, err = x.parseGPSLatitude()
	if err != nil {
		return 0, 0, err
	}

	long, err = x.parseGPSLongitude()
	if err != nil {
		return 0, 0, err
	}

	return lat, long, nil
}

// String returns a pretty text representation of the decoded exif data.
func (x *Exif) String() string {
	var buf bytes.Buffer
	for name, tag := range x.Fields {
		fmt.Fprintf(&buf, "%s: %s\n", name, tag)
	}
	return buf.String()
}

// JpegThumbnail returns the jpeg thumbnail if it exists. If it doesn't exist,
// TagNotPresentError will be returned
func (x *Exif) JpegThumbnail() (int64, int64, error) {
	offset, err := x.Get(models.ThumbJPEGInterchangeFormat)
	if err != nil {
		return 0, 0, err
	}
	start, err := offset.Int(0)
	if err != nil {
		return 0, 0, err
	}

	length, err := x.Get(models.ThumbJPEGInterchangeFormatLength)
	if err != nil {
		return 0, 0, err
	}
	l, err := length.Int(0)
	if err != nil {
		return 0, 0, err
	}

	return int64(start), int64(l), nil
	//return 0, x.Raw[start : start+l], nil
}

// MarshalJSON implements the encoding/json.Marshaler interface providing output of
// all EXIF fields present (names and values).
func (x Exif) MarshalJSON() ([]byte, error) {
	return json.Marshal(x.Fields)
}

type appSec struct {
	marker byte
	data   []byte
}

// newAppSec finds marker in r and returns the corresponding application data
// section.
func newAppSec(marker byte, r io.Reader) (*appSec, error) {
	br := bufio.NewReader(r)
	app := &appSec{marker: marker}
	var dataLen int

	// seek to marker
	for dataLen == 0 {
		if _, err := br.ReadBytes(0xFF); err != nil {
			return nil, err
		}
		c, err := br.ReadByte()
		if err != nil {
			return nil, err
		} else if c != marker {
			continue
		}

		dataLenBytes := make([]byte, 2)
		for k := range dataLenBytes {
			c, err := br.ReadByte()
			if err != nil {
				return nil, err
			}
			dataLenBytes[k] = c
		}
		dataLen = int(binary.BigEndian.Uint16(dataLenBytes)) - 2
	}

	// read section data
	nread := 0
	for nread < dataLen {
		s := make([]byte, dataLen-nread)
		n, err := br.Read(s)
		nread += n
		if err != nil && nread < dataLen {
			return nil, err
		}
		app.data = append(app.data, s[:n]...)
	}
	return app, nil
}

// exifReader returns a reader on this appSec with the read cursor advanced to
// the start of the exif's tiff encoded portion.
func (app *appSec) exifReader() (*bytes.Reader, error) {
	if len(app.data) < 6 {
		return nil, errors.New("exif: failed to find exif intro marker")
	}

	// read/check for exif special mark
	exif := app.data[:6]
	if !bytes.Equal(exif, append([]byte("Exif"), 0x00, 0x00)) {
		return nil, errors.New("exif: failed to find exif intro marker")
	}
	return bytes.NewReader(app.data[6:]), nil
}

// DecodeWithParseHeader parses EXIF data from r (a TIFF, JPEG, HEIF, or rawExif)
// and returns a queryable Exif object. After the EXIF data section is
// called and the TIFF structure is decoded, each registered parser is
// called (in order of registration). If one parser returns an error,
// decoding terminates and the remaining parsers are not called.
//
// Limits the length of Exif data read to increase speed from large files.
// The limit is defined in ExifLengthCutoff.
func DecodeWithParseHeader(r io.Reader) (x *Exif, err error) {
	return DecodeWithParseHeaderAndOptions(r, nil)
}

// DecodeWithParseHeaderAndOptions parses EXIF data with configurable options.
func DecodeWithParseHeaderAndOptions(r io.Reader, opts *DecodeOptions) (x *Exif, err error) {
	if opts == nil {
		opts = &DecodeOptions{}
	}
	if opts.MaxExifSize <= 0 {
		opts.MaxExifSize = exifLengthCutoff
	}
	defer func() {
		if state := recover(); state != nil {
			err = fmt.Errorf("Exif Error: %v", err)
		}
	}()
	r2 := io.LimitReader(r, int64(opts.MaxExifSize))
	data, err := io.ReadAll(r2)
	if err != nil {
		return nil, fmt.Errorf("failed to read EXIF data: %w", err)
	}

	foundAt := -1
	for i := 0; i < len(data); i++ {
		if err = checkExifHeader(data[i:]); err == nil {
			foundAt = i
			break
		}
	}
	if err != nil {
		return
	}

	er := bytes.NewReader(data[foundAt:])
	tif, err := tiff.Decode(er)

	er.Seek(0, 0)
	raw, err := io.ReadAll(er)
	if err != nil {
		return nil, decodeError{cause: err}
	}

	// build an exif structure from the tiff
	x = New(tif, raw, opts)

	for i, p := range parsers {
		if err := p.Parse(x); err != nil {
			if _, ok := err.(tiffErrors); ok {
				return x, err
			}
			// This should never happen, as Parse always returns a tiffError
			// for now, but that could change.
			return x, fmt.Errorf("exif: parser %v failed (%v)", i, err)
		}
	}

	return x, nil
}

// CIPA DC-008-2024 (Exif Version 3.0)
// -> https://www.cipa.jp/std/documents/download_e.html?CIPA_DC-008-2024-E
// Fig 6 p 30: Basic Structure of Uncompressed Data Files
// Table 1 p 31: TIFF Headers
func checkExifHeader(data []byte) error {
	if len(data) < 8 {
		return fmt.Errorf("Invalid EXIF header: too short (length=%d)", len(data))
	}

	byteorder := binary.BigEndian.Uint16(data[0:2])
	var order binary.ByteOrder
	switch byteorder {
	case 0x4d4d: // MM aka Motorola
		order = binary.BigEndian
	case 0x4949: // II aka Intel
		order = binary.LittleEndian
	default:
		return fmt.Errorf("Invalid EXIF header: unrecognized byte order %04x", byteorder)
	}

	fortytwo := order.Uint16(data[2:4])
	if fortytwo != 42 {
		return fmt.Errorf("Invalid EXIF header: got %v, want 42", fortytwo)
	}
	return nil
}

package exif

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/rpajarola/exiftools/models"
	"github.com/rpajarola/exiftools/tiff"
)

// Parser allows the registration of custom parsing and field loading
// in the Decode function.
type Parser interface {
	// Parse should read data from x and insert parsed fields into x via
	// LoadTags.
	Parse(x *Exif) error
}

var parsers []Parser

// RegisterParsers registers one or more parsers to be automatically called
// when decoding EXIF data via the Decode function.
func RegisterParsers(ps ...Parser) {
	parsers = append(parsers, ps...)
}

type parser struct{}

func init() {
	RegisterParsers(&parser{})
}

// Parse is the default exif parser.
// Parse reads data from the tiff data in x and populates the tags
// in x. If parsing a sub-IFD fails, the error is recorded and
// parsing continues with the remaining sub-IFDs.
func (p *parser) Parse(x *Exif) error {
	if len(x.Tiff.Dirs) == 0 {
		return errors.New("invalid exif data")
	}
	keepUnknown := x.opts.KeepUnknownTags
	x.LoadTags(x.Tiff.Dirs[0], models.ExifFields, keepUnknown)

	// thumbnails
	if len(x.Tiff.Dirs) >= 2 {
		x.LoadTags(x.Tiff.Dirs[1], models.ThumbnailFields, keepUnknown)
	}

	te := make(tiffErrors)

	// recurse into exif, gps, and interop sub-IFDs
	if err := x.loadSubDir(models.ExifIFDPointer, models.ExifFields); err != nil {
		te[loadExif] = err.Error()
	}
	if err := x.loadSubDir(models.GPSInfoIFDPointer, models.GpsFields); err != nil {
		te[loadGPS] = err.Error()
	}
	if err := x.loadSubDir(models.InteroperabilityIFDPointer, models.InteropFields); err != nil {
		te[loadInteroperability] = err.Error()
	}
	if len(te) > 0 {
		return te
	}
	return nil
}

func (x *Exif) loadSubDir(ptr models.FieldName, fieldMap map[uint16]models.FieldName) error {
	r := bytes.NewReader(x.Raw)

	tag, err := x.Get(ptr)
	if err != nil {
		return nil
	}
	offset, err := tag.Int64(0)
	if err != nil {
		return nil
	}

	_, err = r.Seek(offset, 0)
	if err != nil {
		return fmt.Errorf("exif: seek to sub-IFD %s failed: %v", ptr, err)
	}
	subDir, _, err := tiff.DecodeDir(r, x.Tiff.Order)
	if err != nil {
		return fmt.Errorf("exif: sub-IFD %s decode failed: %v", ptr, err)
	}
	x.LoadTags(subDir, fieldMap, x.opts.KeepUnknownTags)
	return nil
}

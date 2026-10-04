package mknote

import (
	"bytes"
	"encoding/binary"

	"github.com/rpajarola/exiftools/exif"
	"github.com/rpajarola/exiftools/models"
	"github.com/rpajarola/exiftools/tiff"
)

// Olympus is an exif.Parser for Olympus (and OM System) makernote data.
var Olympus = &olympus{}

type olympus struct{}

// Olympus-specific Maker Note fields. The maker note is a small, mostly
// self-contained TIFF-like structure: a vendor preamble ("OLYMPUS\0" or
// "OM SYSTEM\0\0\0", depending on camera age/branding), followed by a
// byte-order mark, a version number (3 or 4) in place of the usual TIFF
// magic, and then a normal 2-byte-count IFD. Unlike Nikon's maker note
// (whose internal offsets are relative to its own byte-order mark), every
// offset here -- including the top-level IFD's own position -- is relative
// to the start of the vendor preamble, i.e. the MakerNote tag's ValOffset.
var (
	OlympusCameraSettingsPtr  models.FieldName = "Olympus.CameraSettingsIFD" // A sub-IFD
	OlympusPreviewImageValid  models.FieldName = "Olympus.PreviewImageValid"
	OlympusPreviewImageStart  models.FieldName = "Olympus.PreviewImageStart"
	OlympusPreviewImageLength models.FieldName = "Olympus.PreviewImageLength"
)

var makerNoteOlympusFields = map[uint16]models.FieldName{
	0x2020: OlympusCameraSettingsPtr,
}

var makerNoteOlympusCameraSettingsFields = map[uint16]models.FieldName{
	0x0100: OlympusPreviewImageValid,
	0x0101: OlympusPreviewImageStart,
	0x0102: OlympusPreviewImageLength,
}

// OlympusPreviewImageTag is the Preview Image Tag for ORF raw files.
var OlympusPreviewImageTag = exif.NewPreviewImageTag(OlympusPreviewImageStart, OlympusPreviewImageLength, models.FieldName("None"))

// olympusPreambles maps each known maker-note preamble to its byte length.
var olympusPreambles = [][]byte{
	[]byte("OLYMPUS\x00"),
	[]byte("OM SYSTEM\x00\x00\x00"),
}

// Parse decodes Olympus/OM System makernote data found in x and adds it to x.
func (*olympus) Parse(x *exif.Exif) error {
	m, err := x.Get(models.MakerNote)
	if err != nil {
		return nil
	}
	var preambleLen int
	for _, p := range olympusPreambles {
		if len(m.Val) >= len(p) && bytes.Equal(m.Val[:len(p)], p) {
			preambleLen = len(p)
			break
		}
	}
	if preambleLen == 0 {
		return nil
	}
	if len(m.Val) < preambleLen+4 {
		return nil
	}

	// The maker note's internal offsets are relative to its own start
	// (m.ValOffset), so decode it without stripping the preamble.
	if len(m.Val) < preambleLen+2 {
		return nil
	}
	var order binary.ByteOrder
	switch string(m.Val[preambleLen : preambleLen+2]) {
	case "II":
		order = binary.LittleEndian
	case "MM":
		order = binary.BigEndian
	default:
		return nil
	}

	mnReader := bytes.NewReader(m.Val)
	if _, err := mnReader.Seek(int64(preambleLen)+4, 0); err != nil {
		return nil
	}
	// A failure anywhere past this point means this file's maker note is
	// unusually shaped (e.g. an old/cheap camera model whose MakerNote
	// tag is declared shorter than the sub-IFD offsets inside it
	// actually reach) rather than that the file itself is unreadable:
	// treat it as "no preview info available" rather than failing the
	// whole Decode(), which would otherwise also take down every other
	// (unrelated) EXIF field this file does have.
	topDir, _, err := tiff.DecodeDir(mnReader, order)
	if err != nil {
		return nil
	}
	x.LoadTags(topDir, makerNoteOlympusFields, false)
	_ = loadSubDir(x, mnReader, OlympusCameraSettingsPtr, makerNoteOlympusCameraSettingsFields)

	// PreviewImageStart is relative to the maker note's own start; adjust
	// it to an absolute file offset so callers can use it directly.
	if previewTag, err := x.Get(OlympusPreviewImageStart); err == nil {
		offset, _ := previewTag.Int64(0)
		previewTag.SetInt(0, offset+int64(m.ValOffset))
		x.Update(OlympusPreviewImageStart, previewTag)
	}

	return nil
}

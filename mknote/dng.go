package mknote

import (
	"github.com/evanoberholster/exiftools/exif"
)

// AdobeDNG is an exif.Parse for DNG subIfds.
var AdobeDNG = &adobeDNG{}
type adobeDNG struct{}

// Parse decodes all Adobe DNG subIFDS found in x and adds it to x.
func (_ *adobeDNG) Parse(x *exif.Exif) error {
	m, err := x.Get(exif.SubIfdsPointer)
	if err != nil {
		return nil
	}
	if !(m.Count > 0) {
		return nil
	}
	subIfds := []map[uint16]exif.FieldName{
		SubIFD0Fields,
		SubIFD1Fields,
		SubIFD2Fields,
	}
	r := bytes.NewReader(x.Raw)
	for i, sub := range subIfds {
		offset, err := m.Int64(i)
		if err != nil {
			return nil
		}
		_, err = r.Seek(offset, 0)
		if err != nil {
			return fmt.Errorf("exif: seek to sub-IFD %s failed: %v", exif.SubIfdsPointer, err)
		}
		subDir, _, err := tiff.DecodeDir(r, x.Tiff.Order)
		if err != nil {
			return fmt.Errorf("exif: sub-IFD %s decode failed: %v", exif.SubIfdsPointer, err)
		}
		x.LoadTags(subDir, sub, false)
	}

	return nil
}

// PreviewImageTags
var (
	SubIfd0PreviewImage = exif.NewPreviewImageTag(subIfdMap("SubIfd0", PreviewImageStart), subIfdMap("SubIfd0", PreviewImageLength), subIfdMap("SubIfd0", Compression))
	SubIfd1PreviewImage = exif.NewPreviewImageTag(subIfdMap("SubIfd1", PreviewImageStart), subIfdMap("SubIfd1", PreviewImageLength), subIfdMap("SubIfd1", Compression))
	SubIfd2JpegFromRaw  = exif.NewPreviewImageTag(subIfdMap("SubIfd2", JpgFromRawStart), subIfdMap("SubIfd2", JpgFromRawLength), subIfdMap("SubIfd2", Compression))
)

// SubIfd-specific fields
var (
	SubfileType        exif.FieldName = "SubfileType"
	ImageWidth         exif.FieldName = "Width"
	ImageLength        exif.FieldName = "ImageLength"
	Compression        exif.FieldName = "Compression"
	JpgFromRawStart    exif.FieldName = "JpgFromRawStart"
	JpgFromRawLength   exif.FieldName = "JpgFromRawLength"
	PreviewImageStart  exif.FieldName = "PreviewImageStart"
	PreviewImageLength exif.FieldName = "PreviewImageLength"
	PreviewColorSpace  exif.FieldName = "PreviewColorSpace"
	PreviewDateTime    exif.FieldName = "PreviewDateTime"
)

func subIfdMap(ifd string, fn exif.FieldName) exif.FieldName {
	return exif.FieldName(ifd + "." + string(fn))
}

var SubIFD0Fields = map[uint16]exif.FieldName{
	0x00fe: subIfdMap("SubIfd0", SubfileType),
	0x0100: subIfdMap("SubIfd0", ImageWidth),
	0x0101: subIfdMap("SubIfd0", ImageLength),
	0x0103: subIfdMap("SubIfd0", Compression),
	0x0111: subIfdMap("SubIfd0", PreviewImageStart),
	0x0117: subIfdMap("SubIfd0", PreviewImageLength),
	0xc71a: subIfdMap("SubIfd0", PreviewColorSpace),
	0xc71b: subIfdMap("SubIfd0", PreviewDateTime),
}

var SubIFD1Fields = map[uint16]exif.FieldName{
	0x00fe: subIfdMap("SubIfd1", SubfileType),
	0x0100: subIfdMap("SubIfd1", ImageWidth),
	0x0101: subIfdMap("SubIfd1", ImageLength),
	0x0103: subIfdMap("SubIfd1", Compression),
	0x0111: subIfdMap("SubIfd1", PreviewImageStart),
	0x0117: subIfdMap("SubIfd1", PreviewImageLength),
	0xc71a: subIfdMap("SubIfd1", PreviewColorSpace),
	0xc71b: subIfdMap("SubIfd1", PreviewDateTime),
}

var SubIFD2Fields = map[uint16]exif.FieldName{
	0x00fe: subIfdMap("SubIfd2", SubfileType),
	0x0100: subIfdMap("SubIfd2", ImageWidth),
	0x0101: subIfdMap("SubIfd2", ImageLength),
	0x0103: subIfdMap("SubIfd2", Compression),
	0x0111: subIfdMap("SubIfd2", JpgFromRawStart),
	0x0117: subIfdMap("SubIfd2", JpgFromRawLength),
	0xc71a: subIfdMap("SubIfd2", PreviewColorSpace),
	0xc71b: subIfdMap("SubIfd2", PreviewDateTime),
}

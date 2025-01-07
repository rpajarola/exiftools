// Package mknote provides makernote parsers that can be used with goexif/exif.
package mknote

import (
	"bytes"
	"fmt"

	"github.com/rpajarola/exiftools/exif"
	"github.com/rpajarola/exiftools/tiff"
)

var (
	// All is a list of all available makernote parsers
	All = []exif.Parser{Apple, Canon, NikonV3, AdobeDNG, Sony}
)

func init() {
	exif.RegisterParsers(All...)
}

func loadSubDir(x *exif.Exif, r *bytes.Reader, ptr exif.FieldName, fieldMap map[uint16]exif.FieldName) error {
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
	x.LoadTags(subDir, fieldMap, false)
	return nil
}

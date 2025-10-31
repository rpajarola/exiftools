package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/rpajarola/exiftools/exif"
	"github.com/rpajarola/exiftools/exif/structured"
	_ "github.com/rpajarola/exiftools/mknote"
	"github.com/rpajarola/exiftools/models"
	"github.com/trimmer-io/go-xmp/xmp"
)

func main() {
	flag.Parse()
	fname := flag.Arg(0)

	meta, err := ExtractMetadata(fname)
	if err != nil {
		fmt.Printf("Failed to extract metadata: %v\n", err)
		return
	}

	outJson, err := json.MarshalIndent(&meta, "", " ")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Parsed Metadata:\n%+v\n", string(outJson))

}

func ExtractMetadata(path string) (*models.ImageMetadata, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("os.Open(%v): %w", path, err)
	}
	defer f.Close() //nolint:errcheck

	// Decode XMP first
	var xmpDoc *xmp.Document
	if _, err := f.Seek(0, io.SeekStart); err == nil {
		decoder := xmp.NewDecoder(f)
		doc := &xmp.Document{}
		if err := decoder.Decode(doc); err == nil {
			xmpDoc = doc
		}
	}
	// Rewind again for EXIF
	if _, err := f.Seek(0, io.SeekStart); err != nil {
		return nil, fmt.Errorf("seek for EXIF: %w", err)
	}

	opts := &exif.DecodeOptions{KeepUnknownTags: true}
	x, err := exif.DecodeWithOptions(f, opts)

	if err == nil {
		return structured.ToStruct(x, structured.WithXMP(xmpDoc))
	}

	// If EXIF failed, but it's just EOF, treat as "no EXIF but maybe XMP"
	if errors.Is(err, io.EOF) {
		meta := &models.ImageMetadata{}
		if opt := structured.WithXMP(xmpDoc); opt != nil {
			_ = opt(meta)
		}
		return meta, nil
	}

	// True failure
	return nil, fmt.Errorf("exif.Decode(%v): %w", path, err)
}

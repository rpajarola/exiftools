package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rpajarola/exiftools/exif"
	_ "github.com/rpajarola/exiftools/mknote"
	"github.com/rpajarola/exiftools/tiff"
)

const testDataDir = "testdata"

func main() {
	flag.Parse()
	fname := flag.Arg(0)

	f, err := os.Open(fname)
	if err != nil {
		fmt.Printf("os.Open(%v): %v\n", fname, err)
		return
	}
	defer f.Close()
	x, err := exif.Decode(f)
	if err != nil {
		fmt.Printf("exif.Decode(%v): %v\n", fname, err)
		return
	}

	x.Walk(walkFunc(func(name exif.FieldName, tag *tiff.Tag) error {
		fmt.Printf("%v: %v\n", name, tag.String())
		return nil
	}))
}

type walkFunc func(exif.FieldName, *tiff.Tag) error

func (f walkFunc) Walk(name exif.FieldName, tag *tiff.Tag) error {
	return f(name, tag)
}

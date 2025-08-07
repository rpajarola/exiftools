package main

//go:generate go run regen.go -- regression_data.go
//go:generate go fmt regression_data.go

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/rpajarola/exiftools/exif"
	_ "github.com/rpajarola/exiftools/mknote"
	"github.com/rpajarola/exiftools/models"
	"github.com/rpajarola/exiftools/tiff"
)

const testDataDir = "../testdata"

func main() {
	flag.Parse()
	fname := flag.Arg(0)

	dst, err := os.Create(fname)
	if err != nil {
		log.Fatal(err)
	}
	defer dst.Close()

	dir, err := os.Open(testDataDir)
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()

	names, err := dir.Readdirnames(0)
	if err != nil {
		log.Fatal(err)
	}
	sort.Strings(names)
	for i, name := range names {
		names[i] = filepath.Join(testDataDir, name)
	}
	makeExpected(names, dst)
}

func makeExpected(files []string, w io.Writer) {
	fmt.Fprintf(w, "package main\n\n")
	fmt.Fprintf(w, "var regressExpected = map[string]map[string]string{\n")

	for _, name := range files {
		if filepath.Ext(name) != ".jpg" && filepath.Ext(name) != ".heic" {
			fmt.Printf("skipping %v\n", name)
			continue
		}

		f, err := os.Open(name)
		if err != nil {
			fmt.Printf("open %v: %v\n", name, err)
			continue
		}

		x, err := exif.DecodeWithParseHeader(f)
		if err != nil {
			fmt.Printf("decode %v: %v\n", name, err)
			f.Close()
			continue
		}

		var items []string
		x.Walk(walkFunc(func(name models.FieldName, tag *tiff.Tag) error {
			items = append(items, fmt.Sprintf("\"%v\": %q,\n", name, tag.String()))
			return nil
		}))
		sort.Strings(items)

		fmt.Fprintf(w, "\"%v\": {\n", filepath.Base(name))
		for _, item := range items {
			fmt.Fprint(w, item)
		}
		fmt.Fprintf(w, "},\n")
		f.Close()
	}
	fmt.Fprintf(w, "}")
}

type walkFunc func(models.FieldName, *tiff.Tag) error

func (f walkFunc) Walk(name models.FieldName, tag *tiff.Tag) error {
	return f(name, tag)
}

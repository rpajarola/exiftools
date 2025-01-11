package main

import (
	"flag"
	"fmt"
	"os"

        _ "github.com/trimmer-io/go-xmp/models"
        "github.com/trimmer-io/go-xmp/xmp"
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
	bb, err := xmp.ScanPackets(f)
	if err != nil {
		fmt.Printf("xmp.ScanPackets(%v): %v\n", fname, err)
		return
	}
	d := &xmp.Document{}
	if err := xmp.Unmarshal(bb[0], d); err != nil {
                fmt.Printf("xmp.Unmarshal: %v", err)
		return
        }

	d.DumpNamespaces()
	p, err := d.ListPaths()
	if err != nil {
		return
	}
	for _, v := range p {
		fmt.Printf("%s = %s\n", v.Path.String(), v.Value)
	}
	//for _, v := range d.Nodes() {
	// 	v.Dump(0)
	//}
}

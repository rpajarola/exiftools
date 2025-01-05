package mknote

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/rpajarola/exiftools/exif"
	"github.com/rpajarola/exiftools/tiff"
)

// Sony is an exif.Parser for sony makernote data.
var Sony = &sony{}

type sony struct{}

// Sony specific fields
var (
	SonyShutterCount         exif.FieldName = "Sony.ShutterCount"
	SonyInternalSerialNumber exif.FieldName = "Sony.InternalSerialNumber"
	Sony0x9050               exif.FieldName = "Sony.0x9050"
)

var makerNoteSonyFields = map[uint16]exif.FieldName{
	0x9050: Sony0x9050,
}

var Sony0x9050Fields = map[uint16]exif.FieldName{
	0x32: SonyShutterCount,
	0x7c: SonyInternalSerialNumber,
}

var Sony0x9050NameToFields = map[exif.FieldName]uint16{
	SonyShutterCount:         0x32,
	SonyInternalSerialNumber: 0x7c,
}

// Parse decodes Sony makernote data found in x and adds it to x.
func (_ *sony) Parse(x *exif.Exif) error {
	m, err := x.Get(exif.Make)
	if err != nil {
		return nil
	}
	mk, err := m.StringVal()
	if mk != "SONY" {
		return nil
	}
	m, err = x.Get(exif.MakerNote)
	if err != nil {
		return nil
	}

	// Sony maker notes are a single IFD directory with no header.
	// Reader offsets need to be w.r.t. the original tiff structure.
	buf := bytes.NewReader(append(make([]byte, m.ValOffset), m.Val...))
	buf.Seek(int64(m.ValOffset), 0)
	mkNotesDir, _, err := tiff.DecodeDir(buf, x.Tiff.Order)

	if err != nil {
		return nil
	}
	x.LoadTags(mkNotesDir, makerNoteSonyFields, true)

	m9050Dir := &tiff.Dir{}

	m9050, err := x.Get(Sony0x9050)
	if err != nil {
		return nil
	}
	descramble0x9050(m9050.Val)

	m9050Dir.Tags = append(m9050Dir.Tags, makeTag(Sony0x9050NameToFields[SonyShutterCount], tiff.DTLong, 1, []byte{m9050.Val[0x32], m9050.Val[0x33], m9050.Val[0x34], 0}))
	m9050Dir.Tags = append(m9050Dir.Tags, makeTag(Sony0x9050NameToFields[SonyInternalSerialNumber], tiff.DTLong, 1, []byte{m9050.Val[0x7c], m9050.Val[0x7d], m9050.Val[0x7e], m9050.Val[0x7f]}))

	x.LoadTags(m9050Dir, Sony0x9050Fields, true)
	return nil
}

func descramble0x9050(val []byte) {
	// Sony0x9050 is a scrambled binary data structure with data in fixed offsets
	// Create unscrambling table: scrambled = (d*d*d) % 249
	// values >=249 are stored as is
	tbl := make([]int, 256)
	for i, _ := range tbl {
		if i < 249 {
			tbl[i*i*i%249] = i
		} else {
			tbl[i] = i
		}
	}
	// Descramble
	for i, v := range val {
		val[i] = byte(tbl[v])
	}
}

func makeTag(id uint16, dt tiff.DataType, count uint32, val []byte) *tiff.Tag {
	tagBytes := []byte{}
	tagBytes = binary.LittleEndian.AppendUint16(tagBytes, uint16(id))
	tagBytes = binary.LittleEndian.AppendUint16(tagBytes, uint16(dt))
	tagBytes = binary.LittleEndian.AppendUint32(tagBytes, count)
	tagBytes = append(tagBytes, val...)

	tagReader := bytes.NewReader(tagBytes)
	if tag, err := tiff.DecodeTag(tagReader, binary.LittleEndian); err == nil {
		fmt.Printf("MakeTag %v %v\n", tag.Id, tag)
		return tag
	} else {
		fmt.Printf("MakeTag error: %v\n", err)
	}
	return nil
}

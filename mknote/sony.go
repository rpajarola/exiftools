package mknote

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"strings"

	"github.com/rpajarola/exiftools/exif"
	"github.com/rpajarola/exiftools/tiff"
)

// Sony is an exif.Parser for sony makernote data.
var Sony = &sony{}

type sony struct{}

// Sony specific fields
var (
	SonyShutterCount          exif.FieldName = "Sony.ShutterCount"
	SonyShutterCount2         exif.FieldName = "Sony.ShutterCount2"
	SonyShutterCount3         exif.FieldName = "Sony.ShutterCount3"
	SonyInternalSerialNumber  exif.FieldName = "Sony.InternalSerialNumber"
	SonyInternalSerialNumber2 exif.FieldName = "Sony.InternalSerialNumber2"

	// only used internally
	sony0x9050 exif.FieldName = "Sony.0x9050"
)

var makerNoteSonyFields = map[uint16]exif.FieldName{
	0x9050: sony0x9050,
}

var Sony0x9050Fields = map[uint16]exif.FieldName{}

type SonyDataType int

const (
	SonyUint24 SonyDataType = iota
	SonyUint32
	SonyHexString
)

type SonyBinaryTag struct {
	fieldName exif.FieldName
	offset    int
	models    []string
	dataType  SonyDataType
	length    int
	ifNull    bool
}

var Sony0x9050BinaryTags = []SonyBinaryTag{
	{
		// Exiftool: Tag9050b and Tag9050c
		fieldName: SonyShutterCount,
		offset:    0x3a,
		dataType:  SonyUint24,
		models:    []string{"ILCA-99M2", "ILCE-1", "ILCE-6100", "ILCE-6300", "ILCE-6400", "ILCE-6500", "ILCE-6600", "ILCE-7C", "ILCE-7M3", "ILCE-7M4", "ILCE-7RM2", "ILCE-7RM3", "ILCE-7RM3A", "ILCE-7RM4", "ILCE-7RM4A", "ILCE-7RM5", "ILCE-7SM2", "ILCE-7SM3", "ILCE-9", "ILCE-9M2", "ILME-FX3", "ZV-E10"},
	},
	{
		// Exiftool: Tag9050b and Tag9050c
		fieldName: SonyShutterCount2,
		offset:    0x50, // overlaps with 0x52
		dataType:  SonyUint24,
		models:    []string{"ILCE-1", "ILCE-6100", "ILCE-6400", "ILCE-6600", "ILCE-7C", "ILCE-7M4", "ILCE-7RM4", "ILCE-7RM4A", "ILCE-7RM5", "ILCE-7SM3", "ILCE-9M2", "ILME-FX3", "ZV-E10"},
	},
	{
		// Exiftool: Tag9050b
		fieldName: SonyShutterCount2,
		offset:    0x52,
		dataType:  SonyUint24,
		models:    []string{"ILCE-7M3", "ILCE-7RM3", "ILCE-7RM3A"},
	},
	{
		// Exiftool: Tag9050b
		// E-Mount: ShutterCount and dateTime
		fieldName: SonyShutterCount2,
		offset:    0x58,
		dataType:  SonyUint24,
		models:    []string{"ILCE-6300", "ILCE-6400", "ILCE-6500", "ILCE-6600", "ILCE-7C", "ILCE-7M3", "ILCE-7RM2", "ILCE-7RM3", "ILCE-7RM3A", "ILCE-7RM4", "ILCE-7RM4A", "ILCE-7SM2", "ILCE-9", "ILCE-9M2", "ILCA-99M2", "ZV-E10"},
	},
	{
		// Exiftool: Tag9050b
		fieldName: SonyShutterCount3,
		offset:    0x019f,
		dataType:  SonyUint32,
		models:    []string{"ILCE-6100", "ILCE-6400", "ILCE-6600", "ILCE-7C", "ILCE-7M3", "ILCE-7RM3", "ILCE-7RM3A", "ILCE-7RM4", "ILCE-7RM4A", "ILCE-9", "ILCE-9M2", "ZV-E10"},
	},
	{
		// Exiftool: Tag9050b
		fieldName: SonyShutterCount3,
		offset:    0x01cb,
		dataType:  SonyUint32,
		models:    []string{"ILCE-7RM2", "ILCE-7SM2"},
	},
	{
		// Exiftool: Tag9050b
		fieldName: SonyShutterCount3,
		offset:    0x01cd,
		dataType:  SonyUint32,
		models:    []string{"ILCE-6300", "ILCE-6500", "ILCA-99M2"},
	},
	{
		// Exiftool: Tag9050d
		fieldName: SonyShutterCount3,
		offset:    0x000a,
		dataType:  SonyUint24,
		models:    []string{"ILCE-6700", "ILCE-7CM2", "ILCE-7CR"},
	},
	{
		// 9050a
		fieldName: SonyShutterCount2,
		offset:    0x4c,
		dataType:  SonyUint24,
		models:    []string{"ILCE-7", "ILCE-7R", "ILCE-7S", "ILCE-7M2", "ILCE-5000", "ILCE-5100", "ILCE-6000", "ILCE-WX1"},
	},
	{
		// 9050a
		// all other models
		fieldName: SonyShutterCount,
		offset:    0x32,
		dataType:  SonyUint24,
	},
	{
		// 9050a
		fieldName: SonyShutterCount3,
		offset:    0x01a0,
		dataType:  SonyUint32,
		models:    []string{"ILCE-5100", "ILCE-QX1", "ILCA-68", "ILCA-77M2"},
	},
	{
		// 9050a
		fieldName: SonyShutterCount3,
		offset:    0x01aa,
		dataType:  SonyUint32,
		models:    []string{"SLT-A58", "SLT-A99", "SLT-A99V", "HV", "NEX-3N", "NEX-5R", "NEX-5T", "NEX-6", "NEX-VG900", "NEX-VG30E", "ILCE-3000", "ILCE-3500", "ILCE-5000"},
	},
	{
		// 9050a
		fieldName: SonyShutterCount3,
		offset:    0x01bd,
		dataType:  SonyUint32,
		models:    []string{"SLT-A37", "SLT-A37V", "SLT-A57", "SLT-A57V", "SLT-A65", "SLT-A65V", "SLT-A77", "SLT-A77V", "Lunar", "NEX-F3", "NEX-5N", "NEX-7", "NEX-VG20E"},
	},

	{
		// 9050a
		fieldName: SonyInternalSerialNumber,
		offset:    0x7c,
		dataType:  SonyUint32,
		models:    []string{"ILCE-", "ILCA-", "Lunar", "NEX", "SLT-", "HV"},
	},
	{
		// 9050c
		fieldName: SonyInternalSerialNumber,
		offset:    0x7c,
		dataType:  SonyHexString,
		length:    6,
		models:    []string{"ILCE-1"},
	},
	{
		// 9050a
		fieldName: SonyInternalSerialNumber2,
		offset:    0xf0,
		dataType:  SonyHexString,
		length:    5,
		models:    []string{"SLT-", "HV", "ILCA-"},
	},
	{
		// 9050d
		fieldName: SonyInternalSerialNumber,
		offset:    0x38,
		dataType:  SonyHexString,
		length:    6,
		models:    []string{"ZV-E10M2"},
	},
	{
		// Tag9050b and Tag9050c
		fieldName: SonyInternalSerialNumber,
		offset:    0x88,
		dataType:  SonyHexString,
		length:    6,
		models:    []string{"ILCA-99M2", "ILCE-1", "ILCE-6100", "ILCE-6300", "ILCE-6400", "ILCE-6500", "ILCE-6600", "ILCE-7C", "ILCE-7M3", "ILCE-7M4", "ILCE-7RM2", "ILCE-7RM3", "ILCE-7RM3A", "ILCE-7RM4", "ILCE-7RM4A", "ILCE-7RM5", "ILCE-7SM2", "ILCE-7SM3", "ILCE-9", "ILCE-9M2", "ILME-FX3", "ZV-E10"},
	},
}

// Parse decodes Sony makernote data found in x and adds it to x.
func (_ *sony) Parse(x *exif.Exif) error {
	m, err := x.Get(exif.Model)
	if err != nil {
		return nil
	}
	model, err := m.StringVal()
	m, err = x.Get(exif.Make)
	if err != nil {
		return nil
	}
	mk, err := m.StringVal()
	if mk != "SONY" && mk != "HASSELBLAD"{
		return nil
	}
	m, err = x.Get(exif.MakerNote)
	if err != nil {
		return nil
	}

	if len(m.Val) < 13 {
                return nil
        }
	var offset int64 = 0
        if bytes.Compare(m.Val[:10], []byte("SONY DSC \000")) == 0 ||
	   bytes.Compare(m.Val[:10], []byte("SONY CAM \000")) == 0 ||
	   bytes.Compare(m.Val[:13], []byte("SONY MOBILE \000")) == 0 || // how is this possible?
	   bytes.Compare(m.Val[:11], []byte("\000\000SONY PIC\000")) == 0 ||
	   bytes.Compare(m.Val[:10], []byte("VHAB     \000")) == 0 {
	   offset = 12
	}

	// Sony maker notes are a single IFD directory with no header.
	// Reader offsets need to be w.r.t. the original tiff structure.
	buf := bytes.NewReader(append(make([]byte, m.ValOffset), m.Val...))
	buf.Seek(int64(m.ValOffset)+offset, 0)
	mkNotesDir, _, err := tiff.DecodeDir(buf, x.Tiff.Order)

	if err != nil {
		return nil
	}
	x.LoadTags(mkNotesDir, makerNoteSonyFields, false)

	m9050Dir := &tiff.Dir{}

	m9050, err := x.Get(sony0x9050)
	if err != nil {
		return nil
	}
	descramble0x9050(m9050.Val)

	for i, bt := range Sony0x9050BinaryTags {
		Sony0x9050Fields[uint16(i)] = bt.fieldName
		if t := makeBinaryTag(&bt, model, m9050.Val, i); t != nil {
			m9050Dir.Tags = append(m9050Dir.Tags, t)
		}
	}

	x.LoadTags(m9050Dir, Sony0x9050Fields, true)
	return nil
}

func descramble0x9050(val []byte) {
	// sony0x9050 is a scrambled binary data structure with data in fixed offsets
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

func makeBinaryTag(bt *SonyBinaryTag, model string, buf []byte, id int) *tiff.Tag {
	if !checkModel(bt, model) {
		return nil
	}
	val, dt := sonyEncode(bt, buf)
	if !bt.ifNull && !checkNotNull(val) {
		return nil
	}
	t := tiff.MakeTag(uint16(id), dt, 1, binary.LittleEndian, val)
	return t
}

func checkModel(bt *SonyBinaryTag, model string) bool {
	for _, mprefix := range bt.models {
		if strings.HasPrefix(model, mprefix) {
			return true
		}
	}
	return false
}

func sonyEncode(bt *SonyBinaryTag, buf []byte) ([]byte, tiff.DataType) {
	var val []byte
	var dt tiff.DataType
	switch bt.dataType {
	case SonyUint24:
		val = make([]byte, 4)
		dt = tiff.DTLong
		if len(buf) >= bt.offset + 3 {
			val = []byte{buf[bt.offset], buf[bt.offset+1], buf[bt.offset+2], 0}
		}
	case SonyUint32:
		val = make([]byte, 4)
		dt = tiff.DTLong
		if len(buf) >= bt.offset + 4 {
			val = []byte{buf[bt.offset], buf[bt.offset+1], buf[bt.offset+2], buf[bt.offset+3]}
		}
	case SonyHexString:
		val = make([]byte, hex.EncodedLen(bt.length)+1)
		dt = tiff.DTAscii
		if len(buf) >= bt.offset +bt.length {
			hex.Encode(val, buf[bt.offset:bt.offset+bt.length])
		}
	}
	return val, dt
}

func checkNotNull(b []byte) bool {
	for _, v := range b {
		if v != 0 {
			return true
		}
	}
	return false
}

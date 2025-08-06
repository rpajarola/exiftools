package exif

//go:generate go run regen_regress.go -- regress_expected_test.go
//go:generate go fmt regress_expected_test.go

import (
	"bytes"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rpajarola/exiftools/tiff"
)

var testDataDir = flag.String("test_data_dir", "../testdata", "Directory where the data files for testing are located")

func TestDecode(t *testing.T) {
	f, err := os.Open(*testDataDir)
	if err != nil {
		t.Fatalf("Could not open sample directory '%s': %v", *testDataDir, err)
	}

	names, err := f.Readdirnames(0)
	if err != nil {
		t.Fatalf("Could not read sample directory '%s': %v", *testDataDir, err)
	}

	for _, name := range names {
		if !strings.HasSuffix(name, ".jpg") {
			t.Logf("skipping non .jpg file %v", name)
			continue
		}
		t.Run(name, func(t *testing.T) {
			f, err := os.Open(filepath.Join(*testDataDir, name))
			if err != nil {
				t.Errorf("open: %v", err)
			}

			x, err := Decode(f)
			if err != nil {
				t.Errorf("decode: %v", err)
				return
			} else if x == nil {
				t.Errorf("no error and yet %v was not decoded", name)
				return
			}

			if err := x.Walk(&walker{name, t}); err != nil {
				t.Errorf("walk: %v", err)
			}
		})
	}
}

func TestDecodeRawEXIF(t *testing.T) {
	rawFile := filepath.Join(*testDataDir, "raw.exif")
	raw, err := os.ReadFile(rawFile)
	if err != nil {
		t.Fatal(err)
	}
	x, err := Decode(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	got := map[string]string{}
	err = x.Walk(walkFunc(func(name FieldName, tag *tiff.Tag) error {
		got[fmt.Sprint(name)] = fmt.Sprint(tag)
		return nil
	}))
	if err != nil {
		t.Fatalf("Walk: %v", err)
	}
	// The main part of this test is that Decode can parse an
	// "Exif\x00\x00" body, but also briefly sanity check the
	// parsed results:
	if len(got) < 42 {
		t.Errorf("Got %d tags. Want at least 42. Got: %v", len(got), err)
	}
	if g, w := got["DateTime"], `"2018:04:03 07:54:33"`; g != w {
		t.Errorf("DateTime value = %q; want %q", g, w)
	}
}

type walkFunc func(FieldName, *tiff.Tag) error

func (f walkFunc) Walk(name FieldName, tag *tiff.Tag) error {
	return f(name, tag)
}

type walker struct {
	picName string
	t       *testing.T
}

func (w *walker) Walk(field FieldName, tag *tiff.Tag) error {
	// this needs to be commented out when regenerating regress expected vals
	pic := regressExpected[w.picName]
	if pic == nil {
		return fmt.Errorf("regression data not found")
	}

	exp, ok := pic[string(field)]
	if !ok {
		w.t.Errorf("   regression data does not have field %v", field)
		return nil
	}

	s := tag.String()
	if tag.Count == 1 && s != "\"\"" {
		s = fmt.Sprintf("[%s]", s)
	}
	got := tag.String()

	if exp != got {
		fmt.Println("s: ", s)
		fmt.Printf("len(s)=%v\n", len(s))
		w.t.Errorf("   field %v bad tag: expected '%s', got '%s'", field, exp, got)
	}
	return nil
}

func TestMarshal(t *testing.T) {
	name := filepath.Join(*testDataDir, "sample1.jpg")
	f, err := os.Open(name)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	defer f.Close()

	x, err := Decode(f)
	if err != nil {
		t.Fatal(err)
	}
	if x == nil {
		t.Fatal("bad err")
	}

	b, err := x.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%s", b)
}

func testSingleParseDegreesString(t *testing.T, s string, w float64) {
	g, err := parseTagDegreesString(s)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(w-g) > 1e-10 {
		t.Errorf("Wrong parsing result %s: Want %.12f, got %.12f", s, w, g)
	}
}

func TestParseTagDegreesString(t *testing.T) {
	// semicolon as decimal mark
	testSingleParseDegreesString(t, "52,00000,50,00000,34,01180", 52.842781055556) // comma as separator
	testSingleParseDegreesString(t, "52,00000;50,00000;34,01180", 52.842781055556) // semicolon as separator

	// point as decimal mark
	testSingleParseDegreesString(t, "14.00000,44.00000,34.01180", 14.742781055556) // comma as separator
	testSingleParseDegreesString(t, "14.00000;44.00000;34.01180", 14.742781055556) // semicolon as separator
	testSingleParseDegreesString(t, "14.00000;44.00000,34.01180", 14.742781055556) // mixed separators

	testSingleParseDegreesString(t, "-008.0,30.0,03.6", -8.501) // leading zeros

	// no decimal places
	testSingleParseDegreesString(t, "-10,15,54", -10.265)
	testSingleParseDegreesString(t, "-10;15;54", -10.265)

	// incorrect mix of comma and point as decimal mark
	s := "-17,00000,15.00000,04.80000"
	if _, err := parseTagDegreesString(s); err == nil {
		t.Error("parseTagDegreesString: false positive for " + s)
	}
}

// Make sure we error out early when a tag had a count of MaxUint32
func TestMaxUint32CountError(t *testing.T) {
	name := filepath.Join(*testDataDir, "corrupt/max_uint32_exif.jpg")
	f, err := os.Open(name)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	defer f.Close()

	_, err = Decode(f)
	if err == nil {
		t.Fatal("no error on bad exif data")
	}
	if !strings.Contains(err.Error(), "invalid Count offset") {
		t.Fatal("wrong error:", err.Error())
	}
}

// Make sure we error out early with tag data sizes larger than the image file
func TestHugeTagError(t *testing.T) {
	name := filepath.Join(*testDataDir, "corrupt/huge_tag_exif.jpg")
	f, err := os.Open(name)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	defer f.Close()

	_, err = DecodeWithParseHeader(f)
	if err == nil {
		t.Fatal("no error on bad exif data")
	}
	//if !strings.Contains(err.Error(), "short read") {
	//	t.Fatal("wrong error:", err.Error())
	//}
}

// Check for a 0-length tag value
func TestZeroLengthTagError(t *testing.T) {
	name := filepath.Join(*testDataDir, "corrupt/infinite_loop_exif.jpg")
	f, err := os.Open(name)
	if err != nil {
		t.Fatalf("%v\n", err)
	}
	defer f.Close()

	_, err = DecodeWithParseHeader(f)
	if err == nil {
		t.Fatal("no error on bad exif data")
	}
	//if !strings.Contains(err.Error(), "zero length tag value") {
	//	t.Fatal("wrong error:", err.Error())
	//}
}

// TestDecodeHEIF tests HEIF/HEIC file EXIF extraction
func TestDecodeHEIF(t *testing.T) {
	heifFile := filepath.Join(*testDataDir, "park.heic")
	f, err := os.Open(heifFile)
	if err != nil {
		t.Skipf("HEIF test file not found: %v", err)
		return
	}
	defer f.Close()

	// Test that we can decode HEIF files using DecodeWithParseHeader
	// (regular Decode() doesn't handle HEIF format detection)
	x, err := DecodeWithParseHeader(f)
	if err != nil {
		t.Fatalf("Failed to decode HEIF file: %v", err)
	}

	if x == nil {
		t.Fatal("Decode returned nil exif data for HEIF file")
	}

	// Test that we can extract some basic EXIF data
	// Try to get some common fields to verify EXIF data was extracted
	if make, err := x.Get(Make); err == nil {
		if makeStr, err := make.StringVal(); err == nil {
			t.Logf("Camera make: %s", makeStr)
		}
	}

	if model, err := x.Get(Model); err == nil {
		if modelStr, err := model.StringVal(); err == nil {
			t.Logf("Camera model: %s", modelStr)
		}
	}

	if datetime, err := x.Get(DateTime); err == nil {
		if dtStr, err := datetime.StringVal(); err == nil {
			t.Logf("DateTime: %s", dtStr)
		}
	}

	// Test image dimensions
	width, height := x.GetImageSize()
	if width > 0 && height > 0 {
		t.Logf("Image dimensions: %dx%d", width, height)
	}

	// Test orientation
	if orientation, err := x.GetOrientation(); err == nil {
		t.Logf("Orientation: %v", orientation)
	}

	// Verify we can walk through all the EXIF tags
	tagCount := 0
	err = x.Walk(walkFunc(func(name FieldName, tag *tiff.Tag) error {
		tagCount++
		return nil
	}))
	if err != nil {
		t.Errorf("Failed to walk HEIF EXIF tags: %v", err)
	}

	if tagCount == 0 {
		t.Error("No EXIF tags found in HEIF file")
	} else {
		t.Logf("Found %d EXIF tags in HEIF file", tagCount)
	}
}

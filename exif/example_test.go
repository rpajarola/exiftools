package exif

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDecodeExample(t *testing.T) {
	// Test basic EXIF decoding functionality
	fname := filepath.Join("../testdata", "sample1.jpg")

	f, err := os.Open(fname)
	if err != nil {
		t.Skipf("Test file not found: %v", err)
	}
	defer f.Close()

	x, err := Decode(f)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}

	// Test camera model retrieval
	if camModel, err := x.Get(Model); err == nil {
		if model, err := camModel.StringVal(); err == nil {
			t.Logf("Camera model: %s", model)
		}
	}

	// Test focal length retrieval
	if focal, err := x.Get(FocalLength); err == nil {
		if numer, denom, err := focal.Rat2(0); err == nil {
			t.Logf("Focal length: %v/%v", numer, denom)
		}
	}

	// Test convenience functions
	if tm, err := x.DateTime(); err == nil {
		t.Logf("Taken: %v", tm)
	}

	if lat, long, err := x.LatLong(); err == nil {
		t.Logf("Coordinates: %v, %v", lat, long)
	}
}

func TestDecodeWithParseHeaderExample(t *testing.T) {
	// Test EXIF decoding with header parsing - finds any JPEG in testdata
	testDataDir := "../testdata"

	f, err := os.Open(testDataDir)
	if err != nil {
		t.Skipf("Test data directory not found: %v", err)
	}
	defer f.Close()

	names, err := f.Readdirnames(0)
	if err != nil {
		t.Fatalf("Could not read test directory: %v", err)
	}

	// Find first JPEG file
	var testFile string
	for _, name := range names {
		if strings.HasSuffix(strings.ToLower(name), ".jpg") {
			testFile = filepath.Join(testDataDir, name)
			break
		}
	}

	if testFile == "" {
		t.Skip("No JPEG test files found")
	}

	testF, err := os.Open(testFile)
	if err != nil {
		t.Fatalf("Could not open test file: %v", err)
	}
	defer testF.Close()

	x, err := DecodeWithParseHeader(testF)
	if err != nil {
		t.Fatalf("DecodeWithParseHeader failed: %v", err)
	}

	// Test same functionality as basic Decode
	if camModel, err := x.Get(Model); err == nil {
		if model, err := camModel.StringVal(); err == nil {
			t.Logf("Camera model: %s", model)
		}
	}

	if focal, err := x.Get(FocalLength); err == nil {
		if numer, denom, err := focal.Rat2(0); err == nil {
			t.Logf("Focal length: %v/%v", numer, denom)
		}
	}

	if tm, err := x.DateTime(); err == nil {
		t.Logf("Taken: %v", tm)
	}

	if lat, long, err := x.LatLong(); err == nil {
		t.Logf("Coordinates: %v, %v", lat, long)
	}
}

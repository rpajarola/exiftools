package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/rpajarola/exiftools/exif"
	_ "github.com/rpajarola/exiftools/mknote"
	"github.com/rpajarola/exiftools/models"
	"github.com/rpajarola/exiftools/tiff"
)

func TestRegression(t *testing.T) {
	// Get all test files
	dir, err := os.Open(testDataDir)
	if err != nil {
		t.Fatalf("Could not open test data directory '%s': %v", testDataDir, err)
	}
	defer dir.Close()

	names, err := dir.Readdirnames(0)
	if err != nil {
		t.Fatalf("Could not read test data directory '%s': %v", testDataDir, err)
	}

	// Test each file that we have expected data for
	for _, name := range names {
		want, ok := regressExpected[name]
		if !ok {
			t.Errorf("No regression data found for test file %s", name)
			continue
		}

		t.Run(name, func(t *testing.T) {
			compareFile(t, name, want)
		})
	}
}

func compareFile(t *testing.T, filename string, want map[string]string) {
	filepath := filepath.Join(testDataDir, filename)

	f, err := os.Open(filepath)
	if err != nil {
		t.Fatalf("Could not open test file '%s': %v", filepath, err)
	}
	defer f.Close()

	got := make(map[string]string)

	// Use DecodeWithParseHeader to handle all file types including HEIC
	x, err := exif.DecodeWithParseHeader(f)
	if err != nil {
		got["ERROR"] = fmt.Sprintf("%v", err)
	} else {

		// Extract values from decoded EXIF
		err = x.Walk(walkFunc(func(name models.FieldName, tag *tiff.Tag) error {
			got[string(name)] = tag.String()
			return nil
		}))
		if err != nil {
			t.Fatalf("Could not walk EXIF tags for '%s': %v", filepath, err)
		}
	}

	// Check for missing fields and wrong values
	for field, wantValue := range want {
		gotValue, exists := got[field]
		if !exists {
			t.Errorf("Field %s missing from decoded EXIF (want: %s)", field, wantValue)
			continue
		}
		if gotValue != wantValue {
			t.Errorf("Field %s value mismatch:\n  got: %s, want: %s", field, gotValue, wantValue)
		}
	}

	// Check for unexpected fields
	for field, gotValue := range got {
		if _, exists := want[field]; !exists {
			t.Errorf("Unexpected field %s found in decoded EXIF: %s", field, gotValue)
		}
	}

	// Log summary
	t.Logf("File %s: got %d, want %d", filename, len(got), len(want))
}

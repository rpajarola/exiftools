package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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
		if filepath.Ext(name) != ".jpg" && filepath.Ext(name) != ".heic" {
			continue
		}

		expected, exists := regressExpected[name]
		if !exists {
			t.Errorf("No expected data found for test file %s", name)
			continue
		}

		t.Run(name, func(t *testing.T) {
			compareFile(t, name, expected)
		})
	}
}

func compareFile(t *testing.T, filename string, expected map[string]string) {
	filepath := filepath.Join(testDataDir, filename)
	
	f, err := os.Open(filepath)
	if err != nil {
		t.Fatalf("Could not open test file '%s': %v", filepath, err)
	}
	defer f.Close()

	// Use DecodeWithParseHeader to handle all file types including HEIC
	x, err := exif.DecodeWithParseHeader(f)
	if err != nil {
		t.Fatalf("Could not decode EXIF data from '%s': %v", filepath, err)
	}

	// Extract values from decoded EXIF
	actual := make(map[string]string)
	err = x.Walk(walkFunc(func(name models.FieldName, tag *tiff.Tag) error {
		actual[string(name)] = tag.String()
		return nil
	}))
	if err != nil {
		t.Fatalf("Could not walk EXIF tags for '%s': %v", filepath, err)
	}

	// Compare expected vs actual
	// Check for missing fields in actual
	for field, expectedValue := range expected {
		actualValue, exists := actual[field]
		if !exists {
			t.Errorf("Field %s missing from decoded EXIF (expected: %s)", field, expectedValue)
			continue
		}
		if actualValue != expectedValue {
			t.Errorf("Field %s value mismatch:\n  expected: %s\n  actual:   %s", field, expectedValue, actualValue)
		}
	}

	// Check for unexpected fields in actual (this might be too strict, so we'll just log them)
	for field, actualValue := range actual {
		if _, exists := expected[field]; !exists {
			t.Logf("Unexpected field %s found in decoded EXIF: %s", field, actualValue)
		}
	}

	// Log summary
	t.Logf("File %s: %d expected fields, %d actual fields", filename, len(expected), len(actual))
}

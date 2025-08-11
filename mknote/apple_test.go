package mknote

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rpajarola/exiftools/exif"
	"github.com/rpajarola/exiftools/models"
	"github.com/rpajarola/exiftools/tiff"
)

var testDataDir = "../testdata"

func TestAppleParseValidMakerNote(t *testing.T) {
	// Test with a known Apple device image
	testFile := filepath.Join(testDataDir, "park.heic")

	f, err := os.Open(testFile)
	if err != nil {
		t.Skipf("os.Open(%v): %v", testFile, err)
	}
	defer f.Close()

	ex, err := exif.DecodeWithParseHeader(f)
	if err != nil {
		t.Fatalf("Failed to decode EXIF: %v", err)
	}

	if _, err := ex.Get(AppleMakerNoteVersion); err != nil {
		t.Error("No Apple MakerNote fields found, but expected from Apple device")
	}
}

func TestAppleParseMakerNote(t *testing.T) {
	tests := []struct {
		name      string
		signature []byte
		corrupted bool
		wantErr   bool
	}{
		{
			name: "No MakerNote",
		},
		{
			name:      "Valid Apple iOS signature",
			signature: []byte("Apple iOS\000\000\001"),
		},
		{
			name:      "Invalid signature - different manufacturer",
			signature: []byte("Canon\000\000\000\000\000\000\001"),
		},
		{
			name:      "Invalid signature - partial Apple",
			signature: []byte("Apple\000\000\000\000\000\000\001"),
		},
		{
			name:      "Invalid signature - wrong version",
			signature: []byte("Apple iOS\000\000\002"),
		},
		{
			name:      "Valid signature, Corrupted IFD data",
			signature: []byte("Apple iOS\000\000\001"),
			corrupted: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex := exif.New(nil, nil, nil)

			if tt.signature != nil {
				// Create MakerNote with test signature plus minimal valid IFD
				makerNoteData := make([]byte, len(tt.signature))
				copy(makerNoteData, tt.signature)

				// Add minimal valid IFD structure

				// padding to offset 14
				makerNoteData = append(makerNoteData, 0x00, 0x00)

				if tt.corrupted {
					// invalid entry count
					makerNoteData = append(makerNoteData, 0xff, 0xff)
				} else {
					// 0 entries
					makerNoteData = append(makerNoteData, 0x00, 0x00)
				}

				// next IFD offset (no next IFD)
				makerNoteData = append(makerNoteData, 0x00, 0x00, 0x00, 0x00)

				makerNoteTag := &tiff.Tag{Val: makerNoteData}
				ex.Fields[models.MakerNote] = makerNoteTag
			}
			err := Apple.Parse(ex)
			if err != nil {
				t.Errorf("Apple.Parse: Unexpected error: %v", err)
			}

			// Note: For invalid signatures, the function returns nil (no error)
			// This is the expected behavior - it just doesn't process the MakerNote
		})
	}
}

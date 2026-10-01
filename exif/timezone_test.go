package exif_test

// This file is package exif_test (not exif) so it can import both exif and
// mknote: mknote imports exif, so exif's own internal test file can't.

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rpajarola/exiftools/exif"
	_ "github.com/rpajarola/exiftools/mknote"
)

// TestTimeZone checks TimeZone() against real camera files with known,
// manufacturer-specific timezone fields, cross-checked against exiftool's
// output for the same files. Covers both supported sources: Canon.TimeInfo
// (a signed offset stored in a TIFF-unsigned-Long tag -- Int() zero-extends
// it, so TimeZone() must reinterpret the bits as signed) and Nikon.WorldTime
// (a packed, untyped 4-byte blob, not a normal tag at all).
func TestTimeZone(t *testing.T) {
	tests := []struct {
		file       string
		wantOffset string
	}{
		{"NIKON D300 4085317 69139.NEF", "+01:00"},
		{"Canon_EOS_Rebel_T6_IMG_0501.CR2", "-04:00"},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			f, err := os.Open(filepath.Join("../testdata", tt.file))
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			x, err := exif.DecodeWithOptions(f, &exif.DecodeOptions{KeepUnknownTags: true})
			if err != nil {
				t.Fatalf("Decode: %v", err)
			}
			tz, err := x.TimeZone()
			if err != nil {
				t.Fatalf("TimeZone: %v", err)
			}
			_, offsetSec := time.Now().In(tz).Zone()
			sign := '+'
			if offsetSec < 0 {
				sign = '-'
				offsetSec = -offsetSec
			}
			got := fmt.Sprintf("%c%02d:%02d", sign, offsetSec/3600, (offsetSec%3600)/60)
			if got != tt.wantOffset {
				t.Errorf("TimeZone() offset = %v, want %v", got, tt.wantOffset)
			}
		})
	}
}

// TestDateTime_timezone checks that DateTime() applies TimeZone()'s result
// (per its own doc comment) rather than silently discarding it, for the same
// two files and against the same exiftool-verified expectations.
func TestDateTime_timezone(t *testing.T) {
	tests := []struct {
		file string
		want string
	}{
		{"NIKON D300 4085317 69139.NEF", "2015-12-20T21:54:56+01:00"},
		{"Canon_EOS_Rebel_T6_IMG_0501.CR2", "2023-07-21T12:43:29-04:00"},
	}
	for _, tt := range tests {
		t.Run(tt.file, func(t *testing.T) {
			f, err := os.Open(filepath.Join("../testdata", tt.file))
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			x, err := exif.DecodeWithOptions(f, &exif.DecodeOptions{KeepUnknownTags: true})
			if err != nil {
				t.Fatalf("Decode: %v", err)
			}
			dt, err := x.DateTime()
			if err != nil {
				t.Fatalf("DateTime: %v", err)
			}
			if got := dt.Format(time.RFC3339); got != tt.want {
				t.Errorf("DateTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

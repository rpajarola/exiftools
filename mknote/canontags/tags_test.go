package canontags_test

import (
	"testing"

	"github.com/rpajarola/exiftools/mknote/canontags"
)

// Test CanonModelID Values
func TestCanonModel(t *testing.T) {
	_, err := canontags.CanonModel(0x805)
	if err != nil {
		t.Fatalf("Could not find CanonModelID %v: %v", 0x805, err)
	}

}

// Test CanonLensType Values
func TestCanonLens(t *testing.T) {
	lens1 := canontags.CanonLens(138)
	if lens1 != "Canon EF 28-80mm f/2.8-4L" {
		t.Fatalf("Test Failed: Lens %v not found", lens1[0])
	}
	lens2 := canontags.CanonLens(1)
	t.Log(lens2)
}

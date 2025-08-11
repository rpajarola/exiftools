package exif

import (
	"testing"
	"time"

	"github.com/rpajarola/exiftools/models"
	"github.com/rpajarola/exiftools/tiff"
)

func TestGetImageSize(t *testing.T) {
	tests := []struct {
		name           string
		setupTags      func(*Exif)
		expectedWidth  uint
		expectedHeight uint
	}{
		{
			name: "ImageWidth and ImageLength tags present",
			setupTags: func(e *Exif) {
				e.Fields[models.ImageWidth] = tiff.MakeIntTag(0x0001, 1920)
				e.Fields[models.ImageLength] = tiff.MakeIntTag(0x0002, 1080)
			},
			expectedWidth:  1920,
			expectedHeight: 1080,
		},
		{
			name: "PixelXDimension and PixelYDimension tags present",
			setupTags: func(e *Exif) {
				e.Fields[models.PixelXDimension] = tiff.MakeIntTag(0x0001, 3840)
				e.Fields[models.PixelYDimension] = tiff.MakeIntTag(0x0001, 2160)
			},
			expectedWidth:  3840,
			expectedHeight: 2160,
		},
		{
			name: "No tags present",
			setupTags: func(e *Exif) {
				// No tags added
			},
			expectedWidth:  0,
			expectedHeight: 0,
		},
		{
			name: "Only width tag present",
			setupTags: func(e *Exif) {
				e.Fields[models.ImageWidth] = tiff.MakeIntTag(0x0001, 1920)
			},
			expectedWidth:  1920,
			expectedHeight: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exif := New(nil, nil, nil)
			tt.setupTags(exif)

			width, height := exif.GetImageSize()

			if width != tt.expectedWidth {
				t.Errorf("expected width %d, got %d", tt.expectedWidth, width)
			}
			if height != tt.expectedHeight {
				t.Errorf("expected height %d, got %d", tt.expectedHeight, height)
			}
		})
	}
}

func TestGetOrientation(t *testing.T) {
	tests := []struct {
		name        string
		setupTags   func(*Exif)
		expected    models.Orientation
		expectError bool
	}{
		{
			name: "Valid orientation value",
			setupTags: func(e *Exif) {
				e.Fields[models.OrientationTag] = tiff.MakeIntTag(0x0001, 3)
			},
			expected:    models.NewOrientation(3),
			expectError: false,
		},
		{
			name: "No orientation tag",
			setupTags: func(e *Exif) {
				// No tags added
			},
			expected:    1,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exif := New(nil, nil, nil)
			tt.setupTags(exif)

			orientation, err := exif.GetOrientation()

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if orientation != tt.expected {
				t.Errorf("expected orientation %v, got %v", tt.expected, orientation)
			}
		})
	}
}

func TestGetFlashMode(t *testing.T) {
	tests := []struct {
		name        string
		setupTags   func(*Exif)
		expected    models.FlashMode
		expectError bool
	}{
		{
			name: "Valid flash value",
			setupTags: func(e *Exif) {
				e.Fields[models.Flash] = tiff.MakeIntTag(0x0001, 1)
			},
			expected:    models.NewFlashMode(1),
			expectError: false,
		},
		{
			name: "No flash tag",
			setupTags: func(e *Exif) {
				// No tags added
			},
			expected:    models.NoFlashFired,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exif := New(nil, nil, nil)
			tt.setupTags(exif)

			flashMode, err := exif.GetFlashMode()

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if flashMode != tt.expected {
				t.Errorf("expected flash mode %v, got %v", tt.expected, flashMode)
			}
		})
	}
}

func TestGetExposureBias(t *testing.T) {
	tests := []struct {
		name        string
		setupTags   func(*Exif)
		expected    models.ExposureBias
		expectError bool
	}{
		{
			name: "Valid exposure bias",
			setupTags: func(e *Exif) {
				e.Fields[models.ExposureBiasValue] = tiff.MakeRationalTag(0x0003, tiff.Rational{-1, 3})
			},
			expected:    models.NewExposureBias(-1, 3),
			expectError: false,
		},
		{
			name: "No exposure bias tag",
			setupTags: func(e *Exif) {
				// No tags added
			},
			expected:    models.NewExposureBias(0, 0),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exif := New(nil, nil, nil)
			tt.setupTags(exif)

			exposureBias, err := exif.GetExposureBias()

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if exposureBias != tt.expected {
				t.Errorf("expected exposure bias %v, got %v", tt.expected, exposureBias)
			}
		})
	}
}

func TestGetAperture(t *testing.T) {
	tests := []struct {
		name        string
		setupTags   func(*Exif)
		expected    float32
		expectError bool
	}{
		{
			name: "Valid aperture value",
			setupTags: func(e *Exif) {
				e.Fields[models.FNumber] = tiff.MakeRationalTag(0x0003, tiff.Rational{28, 10})
			},
			expected:    2.8,
			expectError: false,
		},
		{
			name: "No FNumber tag",
			setupTags: func(e *Exif) {
				// No tags added
			},
			expected:    0.0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exif := New(nil, nil, nil)
			tt.setupTags(exif)

			aperture, err := exif.GetAperture()

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if aperture != tt.expected {
				t.Errorf("expected aperture %f, got %f", tt.expected, aperture)
			}
		})
	}
}

func TestGetISOSpeed(t *testing.T) {
	tests := []struct {
		name        string
		setupTags   func(*Exif)
		expected    int
		expectError bool
	}{
		{
			name: "Valid ISO value",
			setupTags: func(e *Exif) {
				e.Fields[models.ISOSpeedRatings] = tiff.MakeIntTag(0x0001, 800)
			},
			expected:    800,
			expectError: false,
		},
		{
			name: "No ISO tag",
			setupTags: func(e *Exif) {
				// No tags added
			},
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exif := New(nil, nil, nil)
			tt.setupTags(exif)

			iso, err := exif.GetISOSpeed()

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if iso != tt.expected {
				t.Errorf("expected ISO %d, got %d", tt.expected, iso)
			}
		})
	}
}

func TestGetShutterSpeed(t *testing.T) {
	tests := []struct {
		name        string
		setupTags   func(*Exif)
		expected    models.ShutterSpeed
		expectError bool
	}{
		{
			name: "Valid shutter speed",
			setupTags: func(e *Exif) {
				e.Fields[models.ExposureTime] = tiff.MakeRationalTag(0x0003, tiff.Rational{1, 250})
			},
			expected:    models.NewShutterSpeed(1, 250),
			expectError: false,
		},
		{
			name: "No exposure time tag",
			setupTags: func(e *Exif) {
				// No tags added
			},
			expected:    models.NewShutterSpeed(0, 0),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exif := New(nil, nil, nil)
			tt.setupTags(exif)

			shutterSpeed, err := exif.GetShutterSpeed()

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if shutterSpeed != tt.expected {
				t.Errorf("expected shutter speed %v, got %v", tt.expected, shutterSpeed)
			}
		})
	}
}

func TestGetMeteringMode(t *testing.T) {
	tests := []struct {
		name        string
		setupTags   func(*Exif)
		expected    models.MeteringMode
		expectError bool
	}{
		{
			name: "Valid metering mode",
			setupTags: func(e *Exif) {
				e.Fields[models.MeteringModeTag] = tiff.MakeIntTag(0x0001, 3)
			},
			expected:    models.NewMeteringMode(3),
			expectError: false,
		},
		{
			name: "No metering mode tag",
			setupTags: func(e *Exif) {
				// No tags added
			},
			expected:    models.UnknownMeteringMode,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exif := New(nil, nil, nil)
			tt.setupTags(exif)

			meteringMode, err := exif.GetMeteringMode()

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if meteringMode != tt.expected {
				t.Errorf("expected metering mode %v, got %v", tt.expected, meteringMode)
			}
		})
	}
}

func TestGetExposureMode(t *testing.T) {
	tests := []struct {
		name        string
		setupTags   func(*Exif)
		expected    models.ExposureMode
		expectError bool
	}{
		{
			name: "Valid exposure mode",
			setupTags: func(e *Exif) {
				e.Fields[models.ExposureProgram] = tiff.MakeIntTag(0x0001, 2)
			},
			expected:    models.NewExposureMode(2),
			expectError: false,
		},
		{
			name: "No exposure program tag",
			setupTags: func(e *Exif) {
				// No tags added
			},
			expected:    models.UnkownExposureMode,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exif := New(nil, nil, nil)
			tt.setupTags(exif)

			exposureMode, err := exif.GetExposureMode()

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if exposureMode != tt.expected {
				t.Errorf("expected exposure mode %v, got %v", tt.expected, exposureMode)
			}
		})
	}
}

func TestGetString(t *testing.T) {
	tests := []struct {
		name        string
		setupTags   func(*Exif)
		field       models.FieldName
		expected    string
		expectError bool
	}{
		{
			name: "Valid string value",
			setupTags: func(e *Exif) {
				e.Fields[models.Make] = tiff.MakeAsciiTag(0x0002, "  Canon EOS 5D  ")
			},
			field:       models.Make,
			expected:    "Canon EOS 5D",
			expectError: false,
		},
		{
			name: "Tag not present",
			setupTags: func(e *Exif) {
				// No tags added
			},
			field:       models.Make,
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exif := New(nil, nil, nil)
			tt.setupTags(exif)

			result, err := exif.GetString(tt.field)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected string %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGetStrings(t *testing.T) {
	tests := []struct {
		name        string
		setupTags   func(*Exif)
		fields      []models.FieldName
		expected    string
		expectError bool
	}{
		{
			name: "First field present",
			setupTags: func(e *Exif) {
				e.Fields[models.Make] = tiff.MakeAsciiTag(0x0002, "  Test Value  ")
			},
			fields:      []models.FieldName{models.Make, models.Model},
			expected:    "Test Value",
			expectError: false,
		},
		{
			name: "Second field present",
			setupTags: func(e *Exif) {
				e.Fields[models.Model] = tiff.MakeAsciiTag(0x0002, "  Model Value  ")
			},
			fields:      []models.FieldName{models.Make, models.Model},
			expected:    "Model Value",
			expectError: false,
		},
		{
			name: "No fields present",
			setupTags: func(e *Exif) {
				// No tags added
			},
			fields:      []models.FieldName{models.Make, models.Model},
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exif := New(nil, nil, nil)
			tt.setupTags(exif)

			result, err := exif.GetStrings(tt.fields...)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected string %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestGetUints(t *testing.T) {
	tests := []struct {
		name        string
		setupTags   func(*Exif)
		fields      []models.FieldName
		expected    uint
		expectError bool
	}{
		{
			name: "First field present",
			setupTags: func(e *Exif) {
				e.Fields[models.ImageWidth] = tiff.MakeIntTag(0x0001, 1920)
			},
			fields:      []models.FieldName{models.ImageWidth, models.ImageLength},
			expected:    1920,
			expectError: false,
		},
		{
			name: "Second field present",
			setupTags: func(e *Exif) {
				e.Fields[models.ImageLength] = tiff.MakeIntTag(0x0001, 1080)
			},
			fields:      []models.FieldName{models.ImageWidth, models.ImageLength},
			expected:    1080,
			expectError: false,
		},
		{
			name: "No fields present",
			setupTags: func(e *Exif) {
				// No tags added
			},
			fields:      []models.FieldName{models.ImageWidth, models.ImageLength},
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exif := New(nil, nil, nil)
			tt.setupTags(exif)

			result, err := exif.GetUints(tt.fields...)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected uint %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestGPSAltitude(t *testing.T) {
	tests := []struct {
		name        string
		setupTags   func(*Exif)
		expected    float32
		expectError bool
	}{
		{
			name: "Above sea level",
			setupTags: func(e *Exif) {
				e.Fields[models.GPSAltitude] = tiff.MakeRationalTag(0x0003, tiff.Rational{1500, 10})
				e.Fields[models.GPSAltitudeRef] = tiff.MakeIntTag(0x0001, 0)
			},
			expected:    150.0,
			expectError: false,
		},
		{
			name: "Below sea level",
			setupTags: func(e *Exif) {
				e.Fields[models.GPSAltitude] = tiff.MakeRationalTag(0x0003, tiff.Rational{300, 10})
				e.Fields[models.GPSAltitudeRef] = tiff.MakeIntTag(0x0001, 1)
			},
			expected:    -30.0,
			expectError: false,
		},
		{
			name: "No altitude tag",
			setupTags: func(e *Exif) {
				// No tags added
			},
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exif := New(nil, nil, nil)
			tt.setupTags(exif)

			altitude, err := exif.GPSAltitude()

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if altitude != tt.expected {
				t.Errorf("expected altitude %f, got %f", tt.expected, altitude)
			}
		})
	}
}

func TestFocalLength(t *testing.T) {
	tests := []struct {
		name        string
		setupTags   func(*Exif)
		field       models.FieldName
		expected    float32
		expectError bool
	}{
		{
			name: "Rational type focal length",
			setupTags: func(e *Exif) {
				e.Fields[models.FocalLength] = tiff.MakeRationalTag(0x0003, tiff.Rational{85, 1})
			},
			field:       models.FocalLength,
			expected:    85.0,
			expectError: false,
		},
		{
			name: "Zero denominator",
			setupTags: func(e *Exif) {
				e.Fields[models.FocalLength] = tiff.MakeRationalTag(0x0003, tiff.Rational{85, 0})
			},
			field:       models.FocalLength,
			expected:    85.0,
			expectError: false,
		},
		{
			name: "Short type focal length - single value",
			setupTags: func(e *Exif) {
				e.Fields[models.FocalLength] = tiff.MakeShortTag(0x0005, 0, 85) // When a == 0, returns b as float
			},
			field:       models.FocalLength,
			expected:    85.0,
			expectError: false,
		},
		{
			name: "No focal length tag",
			setupTags: func(e *Exif) {
				// No tags added
			},
			field:       models.FocalLength,
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exif := New(nil, nil, nil)
			tt.setupTags(exif)

			focalLength, err := exif.FocalLength(tt.field)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if focalLength != tt.expected {
				t.Errorf("expected focal length %f, got %f", tt.expected, focalLength)
			}
		})
	}
}

func TestGPSTimeStamp(t *testing.T) {
	tests := []struct {
		name        string
		setupTags   func(*Exif)
		expected    time.Time
		expectError bool
	}{
		{
			name: "Valid GPS timestamp with seconds",
			setupTags: func(e *Exif) {
				e.Fields[models.GPSDateStamp] = tiff.MakeAsciiTag(0x0002, "2023:12:25")
				e.Fields[models.GPSTimeStamp] = tiff.MakeGPSTimeTag(0x0004, 14, 30, 45500, 1000)
			},
			expected:    time.Date(2023, 12, 25, 14, 30, 45, 500000000, time.UTC),
			expectError: false,
		},
		{
			name: "No date stamp tag",
			setupTags: func(e *Exif) {
				// No tags added
			},
			expected:    time.Time{},
			expectError: true,
		},
		{
			name: "No time stamp tag",
			setupTags: func(e *Exif) {
				e.Fields[models.GPSDateStamp] = tiff.MakeAsciiTag(0x0002, "2023:12:25")
			},
			expected:    time.Time{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exif := New(nil, nil, nil)
			tt.setupTags(exif)

			timestamp, err := exif.GPSTimeStamp()

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectError && !timestamp.Equal(tt.expected) {
				t.Errorf("expected timestamp %v, got %v", tt.expected, timestamp)
			}
		})
	}
}

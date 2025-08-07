package exif

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/rpajarola/exiftools/models"
	"github.com/rpajarola/exiftools/tiff"
)

func createMockExif() *Exif {
	return &Exif{
		main: make(map[models.FieldName]*tiff.Tag),
	}
}

func createMockIntTag(value int) *tiff.Tag {
	val := make([]byte, 4)
	binary.LittleEndian.PutUint32(val, uint32(value))
	return tiff.MakeTag(0x0001, tiff.DTLong, 1, binary.LittleEndian, val)
}

func createMockShortTag(values ...int) *tiff.Tag {
	val := make([]byte, 2*len(values))
	for i, v := range values {
		binary.LittleEndian.PutUint16(val[i*2:(i+1)*2], uint16(v))
	}
	return tiff.MakeTag(0x0005, tiff.DTShort, uint32(len(values)), binary.LittleEndian, val)
}

func createMockStringTag(value string) *tiff.Tag {
	val := []byte(value + "\000") // null-terminated
	return tiff.MakeTag(0x0002, tiff.DTAscii, uint32(len(val)), binary.LittleEndian, val)
}

func createMockRationalTag(num, den int64) *tiff.Tag {
	if num < 0 || den < 0 {
		// Use signed rational for negative values
		val := make([]byte, 8)
		binary.LittleEndian.PutUint32(val[0:4], uint32(int32(num)))
		binary.LittleEndian.PutUint32(val[4:8], uint32(int32(den)))
		return tiff.MakeTag(0x0003, tiff.DTSRational, 1, binary.LittleEndian, val)
	} else {
		// Use unsigned rational for positive values
		val := make([]byte, 8)
		binary.LittleEndian.PutUint32(val[0:4], uint32(num))
		binary.LittleEndian.PutUint32(val[4:8], uint32(den))
		return tiff.MakeTag(0x0003, tiff.DTRational, 1, binary.LittleEndian, val)
	}
}

func createMockGPSTimeTag(hour, min int64, secNum, secDen int64) *tiff.Tag {
	val := make([]byte, 24) // 3 rationals = 3 * 8 bytes
	// Hour
	binary.LittleEndian.PutUint32(val[0:4], uint32(hour))
	binary.LittleEndian.PutUint32(val[4:8], 1)
	// Minute  
	binary.LittleEndian.PutUint32(val[8:12], uint32(min))
	binary.LittleEndian.PutUint32(val[12:16], 1)
	// Seconds
	binary.LittleEndian.PutUint32(val[16:20], uint32(secNum))
	binary.LittleEndian.PutUint32(val[20:24], uint32(secDen))
	return tiff.MakeTag(0x0004, tiff.DTRational, 3, binary.LittleEndian, val)
}

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
				e.main[models.ImageWidth] = createMockIntTag(1920)
				e.main[models.ImageLength] = createMockIntTag(1080)
			},
			expectedWidth:  1920,
			expectedHeight: 1080,
		},
		{
			name: "PixelXDimension and PixelYDimension tags present",
			setupTags: func(e *Exif) {
				e.main[models.PixelXDimension] = createMockIntTag(3840)
				e.main[models.PixelYDimension] = createMockIntTag(2160)
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
				e.main[models.ImageWidth] = createMockIntTag(1920)
			},
			expectedWidth:  1920,
			expectedHeight: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exif := createMockExif()
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
				e.main[models.OrientationTag] = createMockIntTag(3)
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
			exif := createMockExif()
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
				e.main[models.Flash] = createMockIntTag(1)
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
			exif := createMockExif()
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
				e.main[models.ExposureBiasValue] = createMockRationalTag(-1, 3)
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
			exif := createMockExif()
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
				e.main[models.FNumber] = createMockRationalTag(28, 10)
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
			exif := createMockExif()
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
				e.main[models.ISOSpeedRatings] = createMockIntTag(800)
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
			exif := createMockExif()
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
				e.main[models.ExposureTime] = createMockRationalTag(1, 250)
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
			exif := createMockExif()
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
				e.main[models.MeteringModeTag] = createMockIntTag(3)
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
			exif := createMockExif()
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
				e.main[models.ExposureProgram] = createMockIntTag(2)
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
			exif := createMockExif()
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
				e.main[models.Make] = createMockStringTag("  Canon EOS 5D  ")
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
			exif := createMockExif()
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
				e.main[models.Make] = createMockStringTag("  Test Value  ")
			},
			fields:      []models.FieldName{models.Make, models.Model},
			expected:    "Test Value",
			expectError: false,
		},
		{
			name: "Second field present",
			setupTags: func(e *Exif) {
				e.main[models.Model] = createMockStringTag("  Model Value  ")
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
			exif := createMockExif()
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
				e.main[models.ImageWidth] = createMockIntTag(1920)
			},
			fields:      []models.FieldName{models.ImageWidth, models.ImageLength},
			expected:    1920,
			expectError: false,
		},
		{
			name: "Second field present",
			setupTags: func(e *Exif) {
				e.main[models.ImageLength] = createMockIntTag(1080)
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
			exif := createMockExif()
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
				e.main[models.GPSAltitude] = createMockRationalTag(1500, 10)
				e.main[models.GPSAltitudeRef] = createMockIntTag(0)
			},
			expected:    150.0,
			expectError: false,
		},
		{
			name: "Below sea level",
			setupTags: func(e *Exif) {
				e.main[models.GPSAltitude] = createMockRationalTag(300, 10)
				e.main[models.GPSAltitudeRef] = createMockIntTag(1)
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
			exif := createMockExif()
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
				e.main[models.FocalLength] = createMockRationalTag(85, 1)
			},
			field:       models.FocalLength,
			expected:    85.0,
			expectError: false,
		},
		{
			name: "Zero denominator",
			setupTags: func(e *Exif) {
				e.main[models.FocalLength] = createMockRationalTag(85, 0)
			},
			field:       models.FocalLength,
			expected:    85.0,
			expectError: false,
		},
		{
			name: "Short type focal length - single value",
			setupTags: func(e *Exif) {
				e.main[models.FocalLength] = createMockShortTag(0, 85) // When a == 0, returns b as float
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
			exif := createMockExif()
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
				e.main[models.GPSDateStamp] = createMockStringTag("2023:12:25")
				e.main[models.GPSTimeStamp] = createMockGPSTimeTag(14, 30, 45500, 1000)
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
				e.main[models.GPSDateStamp] = createMockStringTag("2023:12:25")
			},
			expected:    time.Time{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exif := createMockExif()
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

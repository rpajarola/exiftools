package structured

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rpajarola/exiftools/exif"
	"github.com/rpajarola/exiftools/models"
	"github.com/rpajarola/exiftools/tiff"
)

func rationalToFloat64(tag *tiff.Tag) (*float64, error) {
	num, den, err := tag.Rat2(0)
	if err != nil || den == 0 {
		return nil, err
	}
	val := float64(num) / float64(den)
	return &val, nil
}

// tagString extracts the value of a given EXIF tag as a string, handling both
// textual and numeric formats.
//
// It supports string tags (e.g. "Make", "Model") as well as numeric enumerations
// (e.g. "ExposureProgram", "Contrast", "WhiteBalance"). Known EXIF enums are
// automatically mapped to human-readable labels via models.LookupEnum() and the
// corresponding enum helpers in models (e.g. models.NewFlashMode()).
//
// This helper is tolerant of mixed tag formats and gracefully handles missing or
// unparseable tags by returning nil.
//
// Examples:
//
//	tagString(models.Make, x)                  -> "NIKON"
//	tagString(models.ExposureProgram, x)       -> "Program AE"
//	tagString(models.Flash, x)                 -> "Off, Did not fire"
//	tagString(models.WhiteBalance, x)          -> "Auto"
//	tagString(models.Contrast, x)              -> "Normal"
func tagString(tagID models.FieldName, x *exif.Exif) *string {
	tag, err := x.Get(tagID)
	if err != nil {
		return nil
	}

	switch tag.Format() {
	case tiff.StringVal, tiff.UndefVal:
		// ASCII string tags (e.g., Make, Model, Software)
		if s, err := tag.StringVal(); err == nil {
			str := strings.TrimRight(s, "\x00")
			return &str
		}

		// fallback: raw bytes
		if tag.Val != nil {
			str := strings.TrimRight(string(tag.Val), "\x00 ")
			return &str
		}

	default:
		// Numeric, rational, or enum tags
		if val, err := tag.Int(0); err == nil {
			var s string
			switch tagID {
			case models.OrientationTag:
				s = models.NewOrientation(val).String()
			case models.MeteringModeTag:
				s = models.NewMeteringMode(val).String()
			case models.ExposureModeTag:
				s = models.NewExposureMode(val).String()
			case models.Flash:
				s = models.NewFlashMode(val).String()
			default:
				// Generic enumerations (contrast, saturation, etc.)
				s = models.LookupEnum(tagID, val)
			}

			if s != "" {
				return &s
			}

			// fallback to numeric string if no mapping found
			str := strconv.Itoa(val)
			return &str
		}

		// Rational-valued tags (e.g. ExposureTime)
		if f, err := rationalToFloat64(tag); err == nil && f != nil {
			s := fmt.Sprintf("%v", *f)
			return &s
		}
	}

	return nil
}

func tagInt(tagID models.FieldName, x *exif.Exif) *uint16 {
	if tag, err := x.Get(tagID); err == nil {
		if val, err := tag.Int(0); err == nil {
			v := uint16(val)
			return &v
		}
	}
	return nil
}

func tagRat(tagID models.FieldName, x *exif.Exif) *float64 {
	if tag, err := x.Get(tagID); err == nil {
		if val, err := rationalToFloat64(tag); err == nil {
			return val
		}
	}
	return nil
}

func tagTime(tagID models.FieldName, x *exif.Exif) *time.Time {
	if tag, err := x.Get(tagID); err == nil {
		if str, err := tag.StringVal(); err == nil {
			if t, err := time.Parse("2006:01:02 15:04:05", str); err == nil {
				return &t
			}
		}
	}
	return nil
}

func tagUint32(tagID models.FieldName, x *exif.Exif) *uint32 {
	if tag, err := x.Get(tagID); err == nil {
		if v, err := tag.Int(0); err == nil {
			u := uint32(v)
			return &u
		}
	}
	return nil
}

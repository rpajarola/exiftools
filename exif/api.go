package exif

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/rpajarola/exiftools/models"
	"github.com/rpajarola/exiftools/tiff"
)

// GetImageSize - Get Image Width and Height
func (x *Exif) GetImageSize() (uint, uint) {
	var err error
	var w, l *tiff.Tag

	// Get Width Tag, if error check Pixel Dimension Tag
	w, err = x.Get(models.ImageWidth)
	if err != nil {
		w, err = x.Get(models.PixelXDimension)
		if err != nil {
			return 0, 0
		}
		l, _ = x.Get(models.PixelYDimension)

	} else {
		l, _ = x.Get(models.ImageLength)
	}
	width, _ := w.Int(0)
	length, _ := l.Int(0)
	return uint(width), uint(length)
}

// GetOrientation - Get Image Orientation from Exif
// (Orientation) Tag
func (x *Exif) GetOrientation() (models.Orientation, error) {
	i, err := x.Get(models.OrientationTag)
	if err != nil {
		return 1, err
	}
	o, err := i.Int(0)
	return models.NewOrientation(o), err
}

// GetFlashMode - Get Flash mode from Exif
// (Flash) Tag
func (x *Exif) GetFlashMode() (models.FlashMode, error) {
	tag, err := x.Get(models.Flash)
	if err != nil {
		return models.NoFlashFired, err
	}
	v, err := tag.Int(0)
	if err != nil {
		return models.NoFlashFired, err
	}
	return models.NewFlashMode(v), nil
}

// GetExposureBias - Get Exposure Bias from Exposure Bias Value Tag
func (x *Exif) GetExposureBias() (models.ExposureBias, error) {
	e, err := x.Get(models.ExposureBiasValue)
	if err != nil {
		return models.NewExposureBias(0, 0), err
	}
	num, denom, err := e.Rat2(0)
	return models.NewExposureBias(num, denom), err
}

// GetAperture - Get Aperture from Exif
// (FNumber) tag
func (x *Exif) GetAperture() (float32, error) {
	a, err := x.Get(models.FNumber)
	if err != nil {
		return 0.0, err
	}
	num, denom, err := a.Rat2(0)
	if err != nil {
		return 0.0, err
	}
	return float32(num) / float32(denom), err
}

// GetISOSpeed - Get ISO from Exif
// (ISOSpeedRatings) tag
func (x *Exif) GetISOSpeed() (int, error) {
	a, err := x.Get(models.ISOSpeedRatings)
	if err != nil {
		return 0, err
	}
	i, err := a.Int(0)
	if err != nil {
		return 0.0, err
	}
	return i, err
}

// GetShutterSpeed - Get ShutterSpeed from ExposureTime Tag
func (x *Exif) GetShutterSpeed() (models.ShutterSpeed, error) {
	a, err := x.Get(models.ExposureTime)
	if err != nil {
		return models.NewShutterSpeed(0, 0), err
	}
	nom, denom, err := a.Rat2(0)
	if err != nil {
		return models.NewShutterSpeed(0, 0), err
	}
	return models.NewShutterSpeed(nom, denom), err
}

// GetMeteringMode - Get Metering Mode from MeteringMode Exif Tag
func (x *Exif) GetMeteringMode() (models.MeteringMode, error) {
	tag, err := x.Get(models.MeteringModeTag)
	if err != nil {
		return models.UnknownMeteringMode, err
	}
	v, err := tag.Int(0)
	if err != nil {
		return models.UnknownMeteringMode, err
	}
	return models.NewMeteringMode(v), nil
}

// GetExposureMode - Get Exposure Mode from ExposureProgram Exif Tag
func (x *Exif) GetExposureMode() (models.ExposureMode, error) {
	tag, err := x.Get(models.ExposureProgram)
	if err != nil {
		return models.UnkownExposureMode, err
	}
	v, err := tag.Int(0)
	if err != nil {
		return models.UnkownExposureMode, err
	}
	return models.NewExposureMode(v), nil
}

// GetString - Convenience function for getting String from tagLabel
func (x *Exif) GetString(tagLabel models.FieldName) (string, error) {
	a, err := x.Get(tagLabel)
	if err != nil {
		return "", err
	}

	s, err := a.StringVal()
	if err != nil {
		return "", err
	}
	// Remove leading and trailing white spaces
	return strings.TrimSpace(s), nil
}

// GetStrings -
func (x *Exif) GetStrings(fields ...models.FieldName) (string, error) {
	var err error
	var ok bool
	var a *tiff.Tag
	for _, field := range fields {
		if a, ok = x.main[field]; ok {
			break
		}
	}
	if a == nil {
		return "", TagNotPresentError(fields[len(fields)-1])
	}
	s, err := a.StringVal()
	if err != nil {
		return "", err
	}
	// Remove leading and trailing white spaces
	return strings.TrimSpace(s), err
}

// GetUints - Get Unsigned Ints from fieldNames
func (x *Exif) GetUints(fields ...models.FieldName) (uint, error) {
	var err error
	var ok bool
	var a *tiff.Tag
	for _, field := range fields {
		if a, ok = x.main[field]; ok {
			break
		}
	}
	if a == nil {
		return 0, TagNotPresentError(fields[len(fields)-1])
	}
	i, err := a.Int(0)
	return uint(i), err
}

// GPSAltitude - Convenience function for getting GPSAltitude
func (x *Exif) GPSAltitude() (float32, error) {
	alt, err := x.Get(models.GPSAltitude)
	if err != nil {
		return 0, err
	}

	altRef, err := x.Get(models.GPSAltitudeRef)
	if err != nil {
		return 0, err
	}

	ref, err := altRef.Int(0)
	if err != nil {
		return 0, fmt.Errorf("cannot parse GPS Altitude: %v", err)
	}

	aN, aD, err := alt.Rat2(0)
	if err != nil {
		return 0, fmt.Errorf("cannot parse GPS Altitude: %v", err)
	}

	a := float32(aN / aD)
	if ref == 1 {
		a = a * -1
	} // Adjust for "Below Sea Level"
	return a, nil
}

func calcTimeHelper(n int64, d int64, _ error) string {
	a := int(n / d)
	if a > 10 {
		return strconv.Itoa(a)
	}
	return fmt.Sprintf("0%d", a)
}

// GPSTimeStamp - Convenience function for getting GPS Timestamp
func (x *Exif) GPSTimeStamp() (time.Time, error) {
	var dt time.Time
	dS, err := x.Get(models.GPSDateStamp)
	if err != nil {
		return dt, err
	}
	tS, err := x.Get(models.GPSTimeStamp)
	if err != nil {
		return dt, err
	}
	exifTimeLayout := "2006:01:02"

	dateStr, _ := dS.StringVal()

	hour := calcTimeHelper(tS.Rat2(0))
	min := calcTimeHelper(tS.Rat2(1))

	secN, secD, err := tS.Rat2(2)
	if err != nil {
		return time.ParseInLocation(exifTimeLayout, dateStr, time.UTC)
	}
	exifTimeLayout = "2006:01:02 15:04:05.999"
	sec := float32(secN) / float32(secD)
	if sec < 10 {
		dateStr = fmt.Sprintf("%v %v:%v:0%.3f", dateStr, hour, min, sec)
	} else {
		dateStr = fmt.Sprintf("%v %v:%v:%.3f", dateStr, hour, min, sec)
	}

	return time.ParseInLocation(exifTimeLayout, dateStr, time.UTC)
}

// FocalLength - Convenience function for getting Lens Focal Length
func (x *Exif) FocalLength(fn models.FieldName) (fl float32, err error) {
	tag, err := x.Get(fn)
	if err != nil {
		err = fmt.Errorf("cannot parse Focal Length: %v", err)
		return
	}

	switch tag.Type {

	case tiff.DTRational:
		num, denom, err := tag.Rat2(0)
		if err != nil {
			return 0, err
		}
		if denom == 0 {
			return float32(num), nil
		}
		return float32(num) / float32(denom), nil

	case tiff.DTShort:
		a, _ := tag.Int(0)
		b, err := tag.Int(1)
		if err != nil {
			return float32(a), nil
		}
		l := len(strconv.Itoa(b))

		if a == 0 {
			return float32(b), nil
		}
		if a == 2 {
			if l == 4 {
				return float32(b) / float32(1000), nil
			}
			if l == 3 {
				return float32(b) / float32(100), nil
			}
		}

	}
	return 0, fmt.Errorf("cannot parse FocalLength")
}

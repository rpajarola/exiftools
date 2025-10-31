package structured

import (
	"github.com/rpajarola/exiftools/exif"
	"github.com/rpajarola/exiftools/models"
	"github.com/trimmer-io/go-xmp/models/dc"
	xmpbase "github.com/trimmer-io/go-xmp/models/xmp_base"
	"github.com/trimmer-io/go-xmp/xmp"
)

// ToStructOption allows callers to enrich the resulting ImageMetadata (e.g., add XMP).
type ToStructOption func(meta *models.ImageMetadata) error

// WithXMP populates ImageMetadata.XMP from an XMP document.
// Safe to use even if doc == nil (no-op).
func WithXMP(doc *xmp.Document) func(meta *models.ImageMetadata) error {
	return func(meta *models.ImageMetadata) error {
		if doc == nil {
			return nil
		}

		xmpData := &models.XMPMetadata{}

		// --- Dublin Core ---
		if dcSchema := dc.FindModel(doc); dcSchema != nil {
			if len(dcSchema.Creator) > 0 {
				xmpData.Creator = &dcSchema.Creator[0]
			}
			if dcSchema.Title != nil {
				s := dcSchema.Title.Default()
				xmpData.Title = &s
			}
			if dcSchema.Description != nil {
				s := dcSchema.Description.Default()
				xmpData.Description = &s
			}
			if len(dcSchema.Subject) > 0 {
				xmpData.Keywords = dcSchema.Subject
			}
			// dc:date is an array — use first, if present
			if len(dcSchema.Date) > 0 {
				t := dcSchema.Date[0]
				xmpData.DCDate = &t
			}
		}

		// --- XMP Base ---
		if xb := xmpbase.FindModel(doc); xb != nil {
			// Rating (int), Create/Modify/Metadata dates, CreatorTool
			if xb.Rating != 0 {
				r := xb.Rating
				xmpData.Rating = &r
			}
			if !xb.CreateDate.IsZero() {
				t := xb.CreateDate
				xmpData.CreateDate = &t
			}
			if !xb.ModifyDate.IsZero() {
				t := xb.ModifyDate
				xmpData.ModifyDate = &t
			}
			if !xb.MetadataDate.IsZero() {
				t := xb.MetadataDate
				xmpData.MetadataDate = &t
			}
			if xb.CreatorTool != "" {
				s := xb.CreatorTool
				xmpData.CreatorTool = &s
			}
		}

		// attach only if we actually found anything
		if (xmpData.Creator != nil) || (xmpData.Title != nil) || (xmpData.Description != nil) ||
			len(xmpData.Keywords) > 0 || xmpData.Rating != nil || xmpData.CreateDate != nil ||
			xmpData.ModifyDate != nil || xmpData.MetadataDate != nil || xmpData.CreatorTool != nil ||
			xmpData.DCDate != nil {
			meta.XMP = xmpData
		}
		return nil
	}
}

func ToStruct(x *exif.Exif, opts ...func(*models.ImageMetadata) error) (*models.ImageMetadata, error) {
	meta := &models.ImageMetadata{
		Exif: &models.ExifMetadata{},
	}

	meta.Exif.Make = tagString(models.Make, x)
	meta.Exif.Model = tagString(models.Model, x)
	meta.Exif.Software = tagString(models.Software, x)
	meta.Exif.Artist = tagString(models.Artist, x)
	meta.Exif.Copyright = tagString(models.Copyright, x)
	meta.Exif.ImageDescription = tagString(models.ImageDescription, x)
	meta.Exif.ExposureTime = tagString(models.ExposureTime, x)
	meta.Exif.ExposureProgram = tagString(models.ExposureProgram, x)
	meta.Exif.MeteringMode = tagString(models.MeteringModeTag, x)
	meta.Exif.Flash = tagString(models.Flash, x)
	meta.Exif.WhiteBalance = tagString(models.WhiteBalance, x)
	meta.Exif.ColorSpace = tagString(models.ColorSpace, x)
	meta.Exif.SceneType = tagString(models.SceneType, x)
	meta.Exif.SensingMethod = tagString(models.SensingMethod, x)
	meta.Exif.ResolutionUnit = tagString(models.ResolutionUnit, x)
	meta.Exif.ExifVersion = tagString(models.ExifVersion, x)
	meta.Exif.ComponentsConfiguration = tagString(models.ComponentsConfiguration, x)
	meta.Exif.LensMake = tagString(models.LensMake, x)
	meta.Exif.LensModel = tagString(models.LensModel, x)
	meta.Exif.LensSerialNumber = tagString(models.LensSerialNumber, x)

	meta.Exif.ISOSpeed = tagInt(models.ISOSpeedRatings, x)

	meta.Exif.FNumber = tagRat(models.FNumber, x)
	meta.Exif.ExposureBias = tagRat(models.ExposureBiasValue, x)
	meta.Exif.FocalLength = tagRat(models.FocalLength, x)
	meta.Exif.XResolution = tagRat(models.XResolution, x)
	meta.Exif.YResolution = tagRat(models.YResolution, x)
	meta.Exif.CompressedBitsPerPixel = tagRat(models.CompressedBitsPerPixel, x)

	if t, err := x.DateTime(); err == nil {
		meta.Exif.DateTimeOriginal = &t
	}
	meta.Exif.DateTimeDigitized = tagTime(models.DateTimeDigitized, x)

	if lat, long, err := x.LatLong(); err == nil {
		meta.Exif.GPSLatitude = &lat
		meta.Exif.GPSLongitude = &long
	}

	// Orientation & YCbCrPositioning (they appear in your exiftool dump)
	meta.Exif.Orientation = tagString(models.OrientationTag, x)
	meta.Exif.YCbCrPositioning = tagString(models.YCbCrPositioning, x)

	// FlashpixVersion
	meta.Exif.FlashpixVersion = tagString(models.FlashpixVersion, x)

	// Picture conditions / exposure cluster
	meta.Exif.BrightnessValue = tagRat(models.BrightnessValue, x)
	meta.Exif.ApertureValue = tagRat(models.ApertureValue, x)
	meta.Exif.MaxApertureValue = tagRat(models.MaxApertureValue, x)
	meta.Exif.ShutterSpeedValue = tagRat(models.ShutterSpeedValue, x)
	meta.Exif.ExposureMode = tagString(models.ExposureModeTag, x)
	meta.Exif.LightSource = tagString(models.LightSource, x)
	meta.Exif.GainControl = tagString(models.GainControl, x)
	meta.Exif.Contrast = tagString(models.Contrast, x)
	meta.Exif.Saturation = tagString(models.Saturation, x)
	meta.Exif.Sharpness = tagString(models.Sharpness, x)
	meta.Exif.SceneCaptureType = tagString(models.SceneCaptureType, x)
	meta.Exif.CustomRendered = tagString(models.CustomRendered, x)
	meta.Exif.DigitalZoomRatio = tagRat(models.DigitalZoomRatio, x)
	meta.Exif.FocalLength35mm = tagRat(models.FocalLengthIn35mmFilm, x)

	// Size / dimensions seen in exiftool output
	meta.Exif.ExifImageWidth = tagUint32(models.PixelXDimension, x)
	meta.Exif.ExifImageHeight = tagUint32(models.PixelYDimension, x)

	// UserComment, FileSource
	meta.Exif.UserComment = tagString(models.UserComment, x)
	meta.Exif.FileSource = tagString(models.FileSource, x)

	// GPS extras (your constants exist in models.fields.go)
	meta.Exif.GPSAltitudeRef = tagString(models.GPSAltitudeRef, x)
	meta.Exif.GPSSatelites = tagString(models.GPSSatelites, x)
	meta.Exif.GPSMapDatum = tagString(models.GPSMapDatum, x)
	meta.Exif.GPSTimeStamp = tagString(models.GPSTimeStamp, x)
	meta.Exif.GPSDateStamp = tagString(models.GPSDateStamp, x)

	// If altitude wasn't set yet, set it (exiftool shows it)
	if meta.Exif.GPSAltitude == nil {
		meta.Exif.GPSAltitude = tagRat(models.GPSAltitude, x)
	}

	// Optional: GPS image direction (present in your exiftool dump as ref + value)
	meta.Exif.GPSImgDirectionRef = tagString(models.GPSImgDirectionRef, x)
	meta.Exif.GPSImgDirection = tagRat(models.GPSImgDirection, x)

	meta.Exif.ExposureMode = tagString(models.ExposureModeTag, x)
	meta.Exif.Contrast = tagString(models.Contrast, x)
	meta.Exif.Saturation = tagString(models.Saturation, x)
	meta.Exif.Sharpness = tagString(models.Sharpness, x)
	meta.Exif.SceneCaptureType = tagString(models.SceneCaptureType, x)
	meta.Exif.CustomRendered = tagString(models.CustomRendered, x)
	meta.Exif.LightSource = tagString(models.LightSource, x)
	meta.Exif.GainControl = tagString(models.GainControl, x)
	meta.Exif.FocalLength35mm = tagRat(models.FocalLengthIn35mmFilm, x)

	meta.Exif.BrightnessValue = tagRat(models.BrightnessValue, x)
	meta.Exif.ApertureValue = tagRat(models.ApertureValue, x)
	meta.Exif.ShutterSpeedValue = tagRat(models.ShutterSpeedValue, x)
	meta.Exif.DigitalZoomRatio = tagRat(models.DigitalZoomRatio, x)

	meta.Exif.FileSource = tagString(models.FileSource, x)
	meta.Exif.SceneType = tagString(models.SceneType, x)
	meta.Exif.InteropIndex = tagString(models.InteroperabilityIndex, x)
	meta.Exif.InteropVersion = tagString(models.InteropVersion, x)

	for _, opt := range opts {
		_ = opt(meta)
	}

	return meta, nil
}

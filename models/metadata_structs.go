package models

import (
	"time"

	xmpbase "github.com/trimmer-io/go-xmp/models/xmp_base"
	"github.com/trimmer-io/go-xmp/xmp"
)

// These are optional semantic helpers used by exif/structured.

type ExifMetadata struct {
	Make              *string    `json:"make,omitempty"`
	Model             *string    `json:"model,omitempty"`
	Software          *string    `json:"software,omitempty"`
	Artist            *string    `json:"artist,omitempty"`
	Copyright         *string    `json:"copyright,omitempty"`
	ImageDescription  *string    `json:"image_description,omitempty"`
	DateTimeOriginal  *time.Time `json:"date_time_original,omitempty"`
	DateTimeDigitized *time.Time `json:"date_time_digitized,omitempty"`

	FNumber          *float64 `json:"f_number,omitempty"`
	ExposureTime     *string  `json:"exposure_time,omitempty"`
	ExposureProgram  *string  `json:"exposure_program,omitempty"`
	ExposureBias     *float64 `json:"exposure_bias,omitempty"`
	MeteringMode     *string  `json:"metering_mode,omitempty"`
	ISOSpeed         *uint16  `json:"iso_speed,omitempty"`
	Flash            *string  `json:"flash,omitempty"`
	FocalLength      *float64 `json:"focal_length,omitempty"`
	LensMake         *string  `json:"lens_make,omitempty"`
	LensModel        *string  `json:"lens_model,omitempty"`
	LensSerialNumber *string  `json:"lens_serial_number,omitempty"`

	WhiteBalance  *string `json:"white_balance,omitempty"`
	ColorSpace    *string `json:"color_space,omitempty"`
	Orientation   *string `json:"orientation,omitempty"`
	SceneType     *string `json:"scene_type,omitempty"`
	SensingMethod *string `json:"sensing_method,omitempty"`

	GPSLatitude        *float64 `json:"gps_latitude,omitempty"`
	GPSLongitude       *float64 `json:"gps_longitude,omitempty"`
	GPSAltitude        *float64 `json:"gps_altitude,omitempty"`
	GPSAltitudeRef     *string  `json:"gps_altitude_ref,omitempty"`
	GPSSatelites       *string  `json:"gps_satelites,omitempty"`
	GPSMapDatum        *string  `json:"gps_map_datum,omitempty"`
	GPSTimeStamp       *string  `json:"gps_time_stamp,omitempty"`
	GPSDateStamp       *string  `json:"gps_date_stamp,omitempty"`
	GPSImgDirectionRef *string  `json:"gps_img_direction_ref,omitempty"`
	GPSImgDirection    *float64 `json:"gps_img_direction,omitempty"`

	XResolution    *float64 `json:"x_resolution,omitempty"`
	YResolution    *float64 `json:"y_resolution,omitempty"`
	ResolutionUnit *string  `json:"resolution_unit,omitempty"`

	ExifVersion             *string  `json:"exif_version,omitempty"`
	FlashpixVersion         *string  `json:"flashpix_version,omitempty"`
	ComponentsConfiguration *string  `json:"components_configuration,omitempty"`
	CompressedBitsPerPixel  *float64 `json:"compressed_bits_per_pixel,omitempty"`

	BrightnessValue   *float64 `json:"brightness_value,omitempty"`
	ApertureValue     *float64 `json:"aperture_value,omitempty"`
	MaxApertureValue  *float64 `json:"max_aperture_value,omitempty"`
	ShutterSpeedValue *float64 `json:"shutter_speed_value,omitempty"`
	ExposureMode      *string  `json:"exposure_mode,omitempty"`
	LightSource       *string  `json:"light_source,omitempty"`
	GainControl       *string  `json:"gain_control,omitempty"`
	Contrast          *string  `json:"contrast,omitempty"`
	Saturation        *string  `json:"saturation,omitempty"`
	Sharpness         *string  `json:"sharpness,omitempty"`
	SceneCaptureType  *string  `json:"scene_capture_type,omitempty"`
	CustomRendered    *string  `json:"custom_rendered,omitempty"`
	DigitalZoomRatio  *float64 `json:"digital_zoom_ratio,omitempty"`
	FocalLength35mm   *float64 `json:"focal_length_35mm,omitempty"`

	ExifImageWidth  *uint32 `json:"exif_image_width,omitempty"`
	ExifImageHeight *uint32 `json:"exif_image_height,omitempty"`

	UserComment      *string `json:"user_comment,omitempty"`
	FileSource       *string `json:"file_source,omitempty"`
	YCbCrPositioning *string `json:"ycbcr_positioning,omitempty"`

	InteropIndex   *string `json:"interop_index,omitempty"`
	InteropVersion *string `json:"interop_version,omitempty"`
}

type XMPMetadata struct {
	Creator     *string  `json:"creator,omitempty"`
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`

	Rating       *xmpbase.Rating `json:"rating,omitempty"`
	CreateDate   *xmp.Date       `json:"create_date,omitempty"`
	ModifyDate   *xmp.Date       `json:"modify_date,omitempty"`
	MetadataDate *xmp.Date       `json:"metadata_date,omitempty"`
	CreatorTool  *xmp.AgentName  `json:"creator_tool,omitempty"`

	DCDate *xmp.Date `json:"dc_date,omitempty"`
}

type ImageMetadata struct {
	Exif *ExifMetadata `json:"EXIF,omitempty"`
	XMP  *XMPMetadata  `json:"XMP,omitempty"`
}

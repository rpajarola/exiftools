package models

import "strings"

// LookupEnum returns a human-readable label for the given EXIF tag and numeric value.
// It first resolves standard EXIF enumerations (exact FieldName matches),
// then attempts vendor MakerNote mappings (Nikon/Canon) by tag name.
// If no mapping exists for the tag/value pair, it returns an empty string.
//
// Examples:
//
//	models.LookupEnum(models.ExposureProgram, 2)  == "Program AE"
//	models.LookupEnum(models.WhiteBalance, 0)     == "Auto"
//	models.LookupEnum("Nikon.ToneComp", 1)        == "Normal"
//	models.LookupEnum("ToneComp", 1)              == "Normal"
func LookupEnum(tag FieldName, val int) string {
	name := string(tag)

	// 1) Standard EXIF enums (exact FieldName constants)
	switch tag {
	case ExposureProgram:
		return ExposureProgramValues[val]
	case SceneCaptureType:
		return SceneCaptureTypeValues[val]
	case Contrast:
		return ContrastValues[val]
	case Saturation:
		return SaturationValues[val]
	case Sharpness:
		return SharpnessValues[val]
	case GainControl:
		return GainControlValues[val]
	case LightSource:
		return LightSourceValues[val]
	case ColorSpace:
		return ColorSpaceValues[val]
	case ResolutionUnit:
		return ResolutionUnitValues[val]
	case YCbCrPositioning:
		return YCbCrPositioningValues[val]
	case WhiteBalance:
		if val == 0 {
			return "Auto"
		}
		return "Manual"
	}

	// 2) Nikon MakerNote enums — accept both "Nikon.X" and bare names.
	switch name {
	case "Nikon.WhiteBalance", "WhiteBalance":
		// Only use Nikon mapping for Nikon prefix; otherwise EXIF WhiteBalance handled above.
		if strings.HasPrefix(name, "Nikon.") {
			return NikonWhiteBalanceValues[val]
		}
	case "Nikon.ToneComp", "ToneComp":
		return NikonToneCompValues[val]
	case "Nikon.ColorMode", "ColorMode":
		return NikonColorModeValues[val]
	case "Nikon.ActiveDLighting", "ActiveDLighting":
		return NikonActiveDLightingValues[val]
	case "Nikon.Sharpening", "Sharpening":
		return NikonSharpeningValues[val]
	case "Nikon.ISOSelection", "ISOSelection":
		return NikonISOSelectionValues[val]
	case "Nikon.ImageStabilization", "ImageStabilization":
		return NikonImageStabilizationValues[val]
	case "Nikon.SceneMode", "SceneMode":
		return NikonSceneModeValues[val]
	case "Nikon.ImageAdjustment", "ImageAdjustment":
		return NikonImageAdjustmentValues[val]
	case "Nikon.AuxiliaryLens", "AuxiliaryLens":
		return NikonAuxiliaryLensValues[val]
	case "Nikon.DistortionControl", "DistortionControl":
		return NikonDistortionControlValues[val]
	}

	// 3) Canon MakerNote enums — match by prefix and choose map by suffix.
	if strings.HasPrefix(name, "Canon.") {
		switch {
		case strings.HasSuffix(name, "ExposureMode"):
			return CanonExposureModeValues[val]
		case strings.HasSuffix(name, "MeteringMode"):
			return CanonMeteringModeValues[val]
		case strings.HasSuffix(name, "FocusMode"):
			return CanonFocusModeValues[val]
		case strings.HasSuffix(name, "ContinuousDrive"):
			return CanonContinuousDriveValues[val]
		}
	}

	// No mapping found
	return ""
}

// Common EXIF enumerations mapped to human-readable strings.
// These mirror EXIF 2.3 and maker-specific conventions used by exiftool.

// ExposureProgramValues defines camera program modes.
var ExposureProgramValues = map[int]string{
	0: "Not Defined",
	1: "Manual",
	2: "Program AE",
	3: "Aperture-priority AE",
	4: "Shutter speed priority AE",
	5: "Creative (Slow speed)",
	6: "Action (High speed)",
	7: "Portrait",
	8: "Landscape",
	9: "Bulb",
}

// SceneCaptureTypeValues defines scene modes.
var SceneCaptureTypeValues = map[int]string{
	0: "Standard",
	1: "Landscape",
	2: "Portrait",
	3: "Night Scene",
}

// ContrastValues defines image contrast modes.
var ContrastValues = map[int]string{
	0: "Normal",
	1: "Soft",
	2: "Hard",
}

// SaturationValues defines image saturation modes.
var SaturationValues = map[int]string{
	0: "Normal",
	1: "Low",
	2: "High",
}

// SharpnessValues defines image sharpness modes.
var SharpnessValues = map[int]string{
	0: "Normal",
	1: "Soft",
	2: "Hard",
}

// GainControlValues defines electronic gain adjustments.
var GainControlValues = map[int]string{
	0: "None",
	1: "Low gain up",
	2: "High gain up",
	3: "Low gain down",
	4: "High gain down",
}

// LightSourceValues defines white light sources.
var LightSourceValues = map[int]string{
	0:   "Unknown",
	1:   "Daylight",
	2:   "Fluorescent",
	3:   "Tungsten (Incandescent)",
	4:   "Flash",
	9:   "Fine Weather",
	10:  "Cloudy",
	11:  "Shade",
	12:  "Daylight Fluorescent",
	13:  "Day White Fluorescent",
	14:  "Cool White Fluorescent",
	15:  "White Fluorescent",
	17:  "Standard Light A",
	18:  "Standard Light B",
	19:  "Standard Light C",
	20:  "D55",
	21:  "D65",
	22:  "D75",
	23:  "D50",
	255: "Other",
}

// ColorSpaceValues defines the EXIF color space.
var ColorSpaceValues = map[int]string{
	1:      "sRGB",
	0xFFFF: "Uncalibrated",
}

// ResolutionUnitValues defines units for X/Y resolution.
var ResolutionUnitValues = map[int]string{
	1: "None",
	2: "inches",
	3: "cm",
}

// YCbCrPositioningValues defines chroma subsampling positioning.
var YCbCrPositioningValues = map[int]string{
	1: "Centered",
	2: "Co-sited",
}

package models

// Nikon-specific makernote values and lookup tables.
// These map integer or small code values found in Nikon MakerNote tags
// to human-readable descriptions matching ExifTool and Nikon EXIF specs.
//
// Reference:
//   - CIPA DC-008-2019 (Exif 2.3)
//   - Nikon MakerNote Format (as decoded by ExifTool 12.76)
//   - https://www.exiftool.org/TagNames/Nikon.html
//
// These tags appear under the "MakerNote" IFD and are model-dependent,
// but most Coolpix and DSLR cameras share the same codes below.

var (
	NikonColorModeValues = map[int]string{
		1: "Color",
		2: "Monochrome",
		3: "Sepia",
		4: "Vivid",
		5: "Neutral",
	}

	NikonToneCompValues = map[int]string{
		0: "Auto",
		1: "Normal",
		2: "Low Contrast",
		3: "High Contrast",
	}

	NikonSaturationValues = map[int]string{
		0: "Normal",
		1: "Low",
		2: "High",
	}

	NikonSharpeningValues = map[int]string{
		0: "Auto",
		1: "Normal",
		2: "Low",
		3: "Medium Low",
		4: "Medium High",
		5: "High",
		6: "Extra High",
	}

	NikonWhiteBalanceValues = map[int]string{
		0:  "Auto",
		1:  "Daylight",
		2:  "Cloudy",
		3:  "Tungsten",
		4:  "Fluorescent",
		5:  "Flash",
		6:  "Shade",
		7:  "Kelvin",
		8:  "Manual (Preset)",
		9:  "Custom",
		10: "Auto (2)",
	}

	NikonActiveDLightingValues = map[int]string{
		0: "Off",
		1: "Low",
		2: "Normal",
		3: "High",
		4: "Extra High",
		5: "Auto",
	}

	NikonImageAdjustmentValues = map[int]string{
		0: "Normal",
		1: "Bright+Contrast",
		2: "Custom",
	}

	NikonISOSelectionValues = map[int]string{
		0: "Manual",
		1: "Auto",
	}

	NikonAuxiliaryLensValues = map[int]string{
		0: "Off",
		1: "Wide Adapter",
		2: "Teleconverter",
		3: "Fisheye",
	}

	NikonImageStabilizationValues = map[int]string{
		0: "Off",
		1: "VR-On",
		2: "VR-Active",
	}

	NikonFocusModeValues = map[int]string{
		0: "Manual",
		1: "AF-S (Single)",
		2: "AF-C (Continuous)",
		3: "AF-F (Full-time)",
	}

	NikonSceneModeValues = map[int]string{
		0:  "Standard",
		1:  "Portrait",
		2:  "Landscape",
		3:  "Sports",
		4:  "Night Portrait",
		5:  "Macro",
		6:  "Party/Indoor",
		7:  "Beach/Snow",
		8:  "Sunset",
		9:  "Dusk/Dawn",
		10: "Night Landscape",
	}

	NikonDistortionControlValues = map[int]string{
		0: "Off",
		1: "On",
	}
)

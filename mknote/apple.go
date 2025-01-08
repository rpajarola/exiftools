package mknote

import (
	"bytes"

	"github.com/rpajarola/exiftools/exif"
	"github.com/rpajarola/exiftools/tiff"
)

// Apple is an exif.Parser for Apple makernote data.
var Apple = &apple{}

type apple struct{}

// Apple-specific Maker Note Sub-Ifd pointers
var ()

// Apple-specific Maker Note fields
var (
	AppleMakerNoteVersion                 exif.FieldName = "Apple.MakerNoteVersion"
	AppleAEMatrix                         exif.FieldName = "Apple.AEMatrix"
	AppleRunTime                          exif.FieldName = "Apple.RunTime"
	AppleAEStable                         exif.FieldName = "Apple.AEStable"
	AppleAETarget                         exif.FieldName = "Apple.AETarget"
	AppleAEAverage                        exif.FieldName = "Apple.AEAverage"
	AppleAFStable                         exif.FieldName = "Apple.AFStable"
	AppleAccelerationVector               exif.FieldName = "Apple.AccelerationVector"
	AppleSISMethod                        exif.FieldName = "Apple.SISMethod"
	AppleHDRImageType                     exif.FieldName = "Apple.HDRImageType "
	AppleBurstUUID                        exif.FieldName = "Apple.BurstUUID"
	AppleFocusDistanceRange               exif.FieldName = "Apple.FocusDistanceRange"
	AppleSphereHealthAverageCurrent       exif.FieldName = "Apple.SphereHealthAverageCurrent"
	AppleOrientation                      exif.FieldName = "Apple.Orientation"
	AppleOISMode                          exif.FieldName = "Apple.OISMode"
	AppleSphereStatus                     exif.FieldName = "Apple.SphereStatus"
	AppleContentIdentifier                exif.FieldName = "Apple.ContentIdentifier"
	AppleQRMOutputType                    exif.FieldName = "Apple.QRMOutputType"
	AppleSphereExternalForceOffset        exif.FieldName = "Apple.SphereExternalForceOffset"
	AppleImageCaptureType                 exif.FieldName = "Apple.ImageCaptureType"
	AppleImageUniqueID                    exif.FieldName = "Apple.ImageUniqueID"
	ApplePhotosOriginatingSignature       exif.FieldName = "Apple.PhotosOriginatingSignature"
	AppleLivePhotoVideoIndex              exif.FieldName = "Apple.PhotosOriginatingSignature"
	ApplePhotosRenderOriginatingSignature exif.FieldName = "Apple.PhotosOriginatingSignature"
	AppleImageProcessingFlags             exif.FieldName = "Apple.PhotosOriginatingSignature"
	AppleQualityHint                      exif.FieldName = "Apple.QualityHint"
	ApplePhotosRenderEffect               exif.FieldName = "Apple.PhotosRenderEffect"
	AppleBracketedCaptureSequenceNumber   exif.FieldName = "Apple.BracketedCaptureSequenceNumber"
	AppleLuminanceNoiseAmplitude          exif.FieldName = "Apple.LuminanceNoiseAmplitude"
	AppleOriginatingAppID                 exif.FieldName = "Apple.OriginatingAppID"
	ApplePhotosAppFeatureFlags            exif.FieldName = "Apple.PhotosAppFeatureFlags"
	AppleImageCaptureRequestID            exif.FieldName = "Apple.ImageCaptureRequestID"
	AppleHDRHeadroom                      exif.FieldName = "Apple.HDRHeadroom"
	AppleAFPerformance                    exif.FieldName = "Apple.AFPerformance"
	AppleSceneFlags                       exif.FieldName = "Apple.SceneFlags"
	AppleSignalToNoiseRatioType           exif.FieldName = "Apple.SignalToNoiseRatioType"
	AppleSignalToNoiseRatio               exif.FieldName = "Apple.SignalToNoiseRatio"
	ApplePhotoIdentifier                  exif.FieldName = "Apple.PhotoIdentifier"
	AppleColorTemperature                 exif.FieldName = "Apple.ColorTemperature"
	AppleCameraType                       exif.FieldName = "Apple.CameraType"
	AppleFocusPosition                    exif.FieldName = "Apple.FocusPosition"
	AppleHDRGain                          exif.FieldName = "Apple.HDRGain"
	AppleAFMeasuredDepth                  exif.FieldName = "Apple.AFMeasuredDepth"
	AppleAFConfidence                     exif.FieldName = "Apple.AFConfidence"
	AppleColorCorrectionMatrix            exif.FieldName = "Apple.ColorCorrectionMatrix"
	AppleGreenGhostMitigationStatus       exif.FieldName = "Apple.GreenGhostMitigationStatus"
	AppleSemanticStyle                    exif.FieldName = "Apple.SemanticStyle"
	AppleSemanticStyleRenderingVer        exif.FieldName = "Apple.SemanticStyleRenderingVer"
	AppleSemanticStylePreset              exif.FieldName = "Apple.SemanticStylePreset"
	AppleApple_0x004e                     exif.FieldName = "Apple.Apple_0x004e"
	AppleApple_0x004f                     exif.FieldName = "Apple.Apple_0x004f"
)

// Apple Maker Notes
var makerNoteAppleFields = map[uint16]exif.FieldName{
	0x0001: AppleMakerNoteVersion,
	0x0002: AppleAEMatrix, // PList
	0x0003: AppleRunTime,  // subdir (PList)
	0x0004: AppleAEStable,
	0x0005: AppleAETarget,
	0x0006: AppleAEAverage,
	0x0007: AppleAFStable,
	0x0008: AppleAccelerationVector, // rational
	0x0009: AppleSISMethod,
	0x000a: AppleHDRImageType, // 2=?, 3=HDR, 4=Original
	0x000b: AppleBurstUUID,    // string
	0x000c: AppleFocusDistanceRange,
	0x000d: AppleSphereHealthAverageCurrent,
	0x000e: AppleOrientation, // 0=landcape?, 4=portrait?
	0x000f: AppleOISMode,
	0x0010: AppleSphereStatus,
	0x0011: AppleContentIdentifier,
	0x0012: AppleQRMOutputType,
	0x0013: AppleSphereExternalForceOffset,
	0x0014: AppleImageCaptureType, // 1=ProRAW, 2=Portrait, 10=Photo, 11=manual focus, 12=scene
	0x0015: AppleImageUniqueID,
	0x0016: ApplePhotosOriginatingSignature,
	0x0017: AppleLivePhotoVideoIndex,
	0x0018: ApplePhotosRenderOriginatingSignature,
	0x0019: AppleImageProcessingFlags,
	0x001a: AppleQualityHint,
	0x001b: ApplePhotosRenderEffect,
	0x001c: AppleBracketedCaptureSequenceNumber, // or Flash?
	0x001d: AppleLuminanceNoiseAmplitude,
	0x001e: AppleOriginatingAppID,
	0x001f: ApplePhotosAppFeatureFlags,
	0x0020: AppleImageCaptureRequestID,
	0x0021: AppleHDRHeadroom,
	0x0023: AppleAFPerformance,
	0x0025: AppleSceneFlags,
	0x0026: AppleSignalToNoiseRatioType,
	0x0027: AppleSignalToNoiseRatio,
	0x002b: ApplePhotoIdentifier,
	0x002d: AppleColorTemperature,
	0x002e: AppleCameraType,
	0x002F: AppleFocusPosition,
	0x0030: AppleHDRGain,
	0x0038: AppleAFMeasuredDepth,
	0x003D: AppleAFConfidence,
	0x003E: AppleColorCorrectionMatrix,
	0x003F: AppleGreenGhostMitigationStatus,
	0x0040: AppleSemanticStyle,
	0x0041: AppleSemanticStyleRenderingVer,
	0x0042: AppleSemanticStylePreset,
	0x004e: AppleApple_0x004e, // PList
	0x004f: AppleApple_0x004f, // PList

}

// Parse decodes all Apple makernote data found in x and adds it to x.
func (_ *apple) Parse(x *exif.Exif) error {
	m, err := x.Get(exif.MakerNote)
	if err != nil {
		return nil
	}
	if len(m.Val) < 12 {
		return nil
	}
	if bytes.Compare(m.Val[:12], []byte("Apple iOS\000\000\001")) != 0 {
		return nil
	}

	// Apple makenotes is a self contained IFD with no tiff marker,
	//  offsets are relative to the start of the of the maker note.
	nReader := bytes.NewReader(m.Val)
	nReader.Seek(14, 0)

	mkNotesDir, _, err := tiff.DecodeDir(nReader, x.Tiff.Order)
	if err != nil {
		return err
	}
	// Parse Apple MakerFields
	x.LoadTags(mkNotesDir, makerNoteAppleFields, false)
	return nil
}

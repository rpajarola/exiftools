package mknote

import (
	"bytes"

	"github.com/rpajarola/exiftools/exif"
	"github.com/rpajarola/exiftools/models"
	"github.com/rpajarola/exiftools/tiff"
)

// Apple is an exif.Parser for Apple makernote data.
var Apple = &apple{}

type apple struct{}

// Apple-specific Maker Note Sub-Ifd pointers
var ()

// Apple-specific Maker Note fields
var (
	AppleMakerNoteVersion                 models.FieldName = "Apple.MakerNoteVersion"
	AppleAEMatrix                         models.FieldName = "Apple.AEMatrix"
	AppleRunTime                          models.FieldName = "Apple.RunTime"
	AppleAEStable                         models.FieldName = "Apple.AEStable"
	AppleAETarget                         models.FieldName = "Apple.AETarget"
	AppleAEAverage                        models.FieldName = "Apple.AEAverage"
	AppleAFStable                         models.FieldName = "Apple.AFStable"
	AppleAccelerationVector               models.FieldName = "Apple.AccelerationVector"
	AppleSISMethod                        models.FieldName = "Apple.SISMethod"
	AppleHDRImageType                     models.FieldName = "Apple.HDRImageType "
	AppleBurstUUID                        models.FieldName = "Apple.BurstUUID"
	AppleFocusDistanceRange               models.FieldName = "Apple.FocusDistanceRange"
	AppleSphereHealthAverageCurrent       models.FieldName = "Apple.SphereHealthAverageCurrent"
	AppleOrientation                      models.FieldName = "Apple.Orientation"
	AppleOISMode                          models.FieldName = "Apple.OISMode"
	AppleSphereStatus                     models.FieldName = "Apple.SphereStatus"
	AppleContentIdentifier                models.FieldName = "Apple.ContentIdentifier"
	AppleQRMOutputType                    models.FieldName = "Apple.QRMOutputType"
	AppleSphereExternalForceOffset        models.FieldName = "Apple.SphereExternalForceOffset"
	AppleImageCaptureType                 models.FieldName = "Apple.ImageCaptureType"
	AppleImageUniqueID                    models.FieldName = "Apple.ImageUniqueID"
	ApplePhotosOriginatingSignature       models.FieldName = "Apple.PhotosOriginatingSignature"
	AppleLivePhotoVideoIndex              models.FieldName = "Apple.PhotosOriginatingSignature"
	ApplePhotosRenderOriginatingSignature models.FieldName = "Apple.PhotosOriginatingSignature"
	AppleImageProcessingFlags             models.FieldName = "Apple.PhotosOriginatingSignature"
	AppleQualityHint                      models.FieldName = "Apple.QualityHint"
	ApplePhotosRenderEffect               models.FieldName = "Apple.PhotosRenderEffect"
	AppleBracketedCaptureSequenceNumber   models.FieldName = "Apple.BracketedCaptureSequenceNumber"
	AppleLuminanceNoiseAmplitude          models.FieldName = "Apple.LuminanceNoiseAmplitude"
	AppleOriginatingAppID                 models.FieldName = "Apple.OriginatingAppID"
	ApplePhotosAppFeatureFlags            models.FieldName = "Apple.PhotosAppFeatureFlags"
	AppleImageCaptureRequestID            models.FieldName = "Apple.ImageCaptureRequestID"
	AppleHDRHeadroom                      models.FieldName = "Apple.HDRHeadroom"
	AppleAFPerformance                    models.FieldName = "Apple.AFPerformance"
	AppleSceneFlags                       models.FieldName = "Apple.SceneFlags"
	AppleSignalToNoiseRatioType           models.FieldName = "Apple.SignalToNoiseRatioType"
	AppleSignalToNoiseRatio               models.FieldName = "Apple.SignalToNoiseRatio"
	ApplePhotoIdentifier                  models.FieldName = "Apple.PhotoIdentifier"
	AppleColorTemperature                 models.FieldName = "Apple.ColorTemperature"
	AppleCameraType                       models.FieldName = "Apple.CameraType"
	AppleFocusPosition                    models.FieldName = "Apple.FocusPosition"
	AppleHDRGain                          models.FieldName = "Apple.HDRGain"
	AppleAFMeasuredDepth                  models.FieldName = "Apple.AFMeasuredDepth"
	AppleAFConfidence                     models.FieldName = "Apple.AFConfidence"
	AppleColorCorrectionMatrix            models.FieldName = "Apple.ColorCorrectionMatrix"
	AppleGreenGhostMitigationStatus       models.FieldName = "Apple.GreenGhostMitigationStatus"
	AppleSemanticStyle                    models.FieldName = "Apple.SemanticStyle"
	AppleSemanticStyleRenderingVer        models.FieldName = "Apple.SemanticStyleRenderingVer"
	AppleSemanticStylePreset              models.FieldName = "Apple.SemanticStylePreset"
	AppleApple_0x004e                     models.FieldName = "Apple.Apple_0x004e"
	AppleApple_0x004f                     models.FieldName = "Apple.Apple_0x004f"
)

// Apple Maker Notes
var makerNoteAppleFields = map[uint16]models.FieldName{
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
func (*apple) Parse(x *exif.Exif) error {
	m, err := x.Get(models.MakerNote)
	if err != nil {
		return nil
	}
	if len(m.Val) < 12 {
		return nil
	}
	if !bytes.Equal(m.Val[:12], []byte("Apple iOS\000\000\001")) {
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

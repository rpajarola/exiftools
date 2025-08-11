package exif

import (
	"errors"
	"fmt"
	"strings"

	"github.com/rpajarola/exiftools/models"
	"github.com/rpajarola/exiftools/tiff"
)

// Exif Errors
var (
	ErrNoExif          = errors.New("no exif data")
	ErrExifHeaderError = errors.New("exif header error")
)

// A decodeError is returned when the image cannot be decoded as a tiff image.
type decodeError struct {
	cause error
}

func (de decodeError) Error() string {
	return fmt.Sprintf("exif: decode failed (%v) ", de.cause.Error())
}

// IsShortReadTagValueError identifies a ErrShortReadTagValue error.
func IsShortReadTagValueError(err error) bool {
	de, ok := err.(decodeError)
	if ok {
		return de.cause == tiff.ErrShortReadTagValue
	}
	return false
}

// A TagNotPresentError is returned when the requested field is not
// present in the EXIF.
type TagNotPresentError models.FieldName

func (tag TagNotPresentError) Error() string {
	return fmt.Sprintf("exif: tag %q is not present", string(tag))
}

// IsTagNotPresentError -
func IsTagNotPresentError(err error) bool {
	_, ok := err.(TagNotPresentError)
	return ok
}

type tiffErrors map[tiffError]string
type tiffError int

func (te tiffErrors) Error() string {
	var allErrors []string
	for k, v := range te {
		allErrors = append(allErrors, fmt.Sprintf("%s: %v\n", stagePrefix[k], v))
	}
	return strings.Join(allErrors, "\n")
}

// IsCriticalError - given the error returned by Decode, reports whether the
// returned *Exif may contain usable information.
func IsCriticalError(err error) bool {
	_, ok := err.(tiffErrors)
	return !ok
}

// IsExifError reports whether the error happened while decoding the EXIF
// sub-IFD.
func IsExifError(err error) bool {
	if te, ok := err.(tiffErrors); ok {
		_, isExif := te[loadExif]
		return isExif
	}
	return false
}

// IsGPSError reports whether the error happened while decoding the GPS sub-IFD.
func IsGPSError(err error) bool {
	if te, ok := err.(tiffErrors); ok {
		_, isGPS := te[loadGPS]
		return isGPS
	}
	return false
}

// IsInteroperabilityError reports whether the error happened while decoding the
// Interoperability sub-IFD.
func IsInteroperabilityError(err error) bool {
	if te, ok := err.(tiffErrors); ok {
		_, isInterop := te[loadInteroperability]
		return isInterop
	}
	return false
}

var stagePrefix = map[tiffError]string{
	loadExif:             "loading EXIF sub-IFD",
	loadGPS:              "loading GPS sub-IFD",
	loadInteroperability: "loading Interoperability sub-IFD",
}

const (
	loadExif tiffError = iota
	loadGPS
	loadInteroperability
)

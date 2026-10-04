package exif

import (
	"sort"

	"github.com/rpajarola/exiftools/models"
)

var exifCompressionValues = map[uint16]string{
	6:     "JPEG (old-style)",
	7:     "JPEG",
	99:    "JPEG",
	34712: "JPEG 2000",
	34892: "Lossy JPEG",
	34934: "JPEG XR",
	34927: "WebP",
	34933: "PNG",
}

// PreviewImageTag -
type PreviewImageTag struct {
	StartTag    models.FieldName
	LengthTag   models.FieldName
	Compression models.FieldName
	Start       int
	Length      int
}

// NewPreviewImageTag -
func NewPreviewImageTag(start models.FieldName, length models.FieldName, compression models.FieldName) PreviewImageTag {
	return PreviewImageTag{start, length, compression, 0, 0}
}

// PreviewImage returns the byte start location and length of the largest
// preview image candidate found. The largest candidate isn't always a real,
// decodable preview (e.g. a RAW format's main StripOffsets/StripByteCounts
// tag pair can be bigger than any actual preview while pointing at raw
// sensor data, not an image), so callers that need a usable preview should
// prefer PreviewCandidates and fall back through the list instead.
func (x Exif) PreviewImage(tags ...PreviewImageTag) (start int64, length int64, err error) {
	candidates := x.PreviewCandidates(tags...)
	if len(candidates) == 0 {
		return 0, 0, nil
	}
	return int64(candidates[0].Start), int64(candidates[0].Length), nil
}

// PreviewCandidates returns every plausible preview-image byte range found
// via tags (plus the built-in generic IFD0 PreviewImage/ThumbnailImage
// tags), sorted largest-first. Each candidate's Start/Length describe a
// byte range that the tag metadata claims holds a preview image, but
// callers should still verify it decodes as one: the metadata alone can't
// distinguish an actual JPEG/PNG preview from e.g. raw sensor data that
// happens to be referenced by a same-shaped tag pair.
func (x Exif) PreviewCandidates(tags ...PreviewImageTag) []PreviewImageTag {
	tags = append(tags,
		NewPreviewImageTag(models.PreviewImageStart, models.PreviewImageLength, models.FieldName("None")),                        // IFD0 PreviewImage
		NewPreviewImageTag(models.ThumbJPEGInterchangeFormat, models.ThumbJPEGInterchangeFormatLength, models.FieldName("None")), // IFD0 ThumbnailImage
	)
	var res []PreviewImageTag
	for _, tag := range tags {
		// If Preview Image is of type JPEG, PNG, WEBP else continue
		if tag.Compression != models.FieldName("None") {
			compression, err := x.Get(tag.Compression)
			if err == nil {
				c, err := compression.Int(0)
				if err != nil {
					continue
				}
				_, ok := exifCompressionValues[uint16(c)]
				if !ok {
					continue
				}
			}
		}
		offset, err := x.Get(tag.StartTag)
		if err != nil {
			continue
		}
		tag.Start, err = offset.Int(0)
		if err != nil {
			continue
		}
		length, err := x.Get(tag.LengthTag)
		if err != nil {
			continue
		}
		tag.Length, err = length.Int(0)
		if err != nil {
			continue
		}
		if tag.Length <= 0 {
			continue
		}
		res = append(res, tag)
	}

	sort.Slice(res, func(i, j int) bool { return res[i].Length > res[j].Length })
	return res
}

// PreviewBlobCandidates returns the raw value bytes of each present tag in
// names, sorted largest-first. Unlike a PreviewImageTag start/length pair,
// a "blob" tag (e.g. Panasonic RW2's JpgFromRaw) holds or references the
// preview image data directly as its own value, with no separate length
// tag to pair it with.
func (x Exif) PreviewBlobCandidates(names ...models.FieldName) [][]byte {
	var blobs [][]byte
	for _, name := range names {
		tag, err := x.Get(name)
		if err != nil || len(tag.Val) == 0 {
			continue
		}
		blobs = append(blobs, tag.Val)
	}
	sort.Slice(blobs, func(i, j int) bool { return len(blobs[i]) > len(blobs[j]) })
	return blobs
}

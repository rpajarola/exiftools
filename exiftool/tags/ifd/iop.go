package ifd

import (
	"github.com/rpajarola/exiftools/exiftool/exif"
)

// IfdIop Name and TagID
const (
	IfdIop              = "Iop"
	IfdIopID exif.TagID = 0xA005
)

// IopPath is the IFD/Iop Ifd Path
var (
	IopPath = IfdPath{IfdRootID}
)

// IopIfd is the IFD/Iop IFD for Interoperability Information
var IopIfd = IfdItem{IopPath, IfdIopID, IfdIop}

# Known issues

`go vet ./...` and `staticcheck ./...` are both clean. `go test ./...` passes across
every package. Everything found by the original correctness review (uint32 overflow in
the tag-size bounds check, IFD cycle detection missing longer cycles, `tag.Rat2`
panicking on an empty rational tag, the dead-code/style staticcheck findings, and
`TestDecode`'s no-EXIF assumption) is fixed as of this file's git history.

## Pre-existing TODOs already in the source

Feature gaps, not bugs — left as-is rather than fixed opportunistically, since they're
more involved than a one-line correctness fix and the right behavior isn't obvious
without more context on intent:

- `exif/exif.go:371,382` — timezone parsing gaps (GPS time, Nikon WorldTime)
- `mknote/nikon.go:30` — "fix regression test"
- `tiff/tag.go:178` — tag values pointing outside the current IFD block aren't handled;
  per the comment there, this is common for thumbnails and other large blobs

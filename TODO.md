# Known issues

`go vet ./...` and `staticcheck ./...` are both clean. `go test ./...` passes across
every package. Everything found by the original correctness review, every staticcheck
finding, `TestDecode`'s no-EXIF assumption, and all but one of the pre-existing
source TODOs are fixed as of this file's git history (DateTime/TimeZone's Canon sign
bug and missing Nikon WorldTime support, and NikonISOInfo's regression-test breakage,
which no longer reproduces).

## `tiff/tag.go:178` — tag values pointing outside the current IFD block

When a tag's declared value is larger than 4 bytes, `DecodeTag` reads it via a
`SectionReader` over `r` at the tag's absolute `ValOffset`. If that read comes up short
(`n != valLen`), the value is silently dropped (`t.Val` stays empty) rather than erroring,
with a comment noting this is "common for thumbnails and other large blobs."

Investigated by instrumenting this exact branch and running it across the full shared
test corpus (both repos' `testdata/`, including real Sony/Canon/Nikon/Apple/Samsung
files with legitimate embedded thumbnails) — it never fired. Not a crash or correctness
risk as it stands (the short-read path returns `t, nil`, which `DecodeDir`'s tag loop
accepts and continues past; switching it to return the already-declared
`ErrShortReadTagValue` would need `DecodeDir` to also treat that error as skip-not-abort,
same as it already does for `errUnhandledTagType`, or every other tag in the same IFD
would be lost along with it). Left as-is: no reproducible case to verify a fix against,
and the real-world thumbnail/preview case this comment worries about is already handled
by a different, working mechanism (`Exif.Raw` + `PreviewImage()`/`JpegThumbnail()`, which
read the full captured buffer directly rather than going through `Tag.Val`).

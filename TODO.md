# Known issues

Found via a correctness-focused review of the core parsing packages (`tiff/`, `exif/`,
`mknote/`), plus `go vet`/`staticcheck`.

## Unrecovered panic on malformed-but-plausible GPS tag data

Reachable through a normal public API call, not deep internals.

### `exif/exif.go:458` (`parse3Rat2`) calls `tag.Rat2(i)` before checking `tag.Count`

The bounds-ish check (`if tag.Count < uint32(i+2) { break }`) runs *after* `tag.Rat2(i)`
is already called and its result used. For `i == 0` with `tag.Count == 0` (a GPS tag
present but declaring zero rational values), `tag.Rat2(0)` does `t.ratVals[0][0]` on an
empty slice and panics. Unlike `Int()` (which has `defer recover()`) or `Int64()` (which
explicitly bounds-checks), `Rat2()` (`tiff/tag.go:451-456`) has neither. Reachable via the
public `LatLong()` / GPS-parsing path.

## Lower-severity findings (staticcheck / go vet)

`go vet ./...` is clean. `staticcheck ./...` found only style/dead-code items, listed
here rather than fixed since none affect correctness:

- `dump_exif/main.go:14`, `dump_xmp/main.go:12` — unused `const testDataDir`
- `example/main.go:155` — unused `func colorJSON`
- `mknote/canon.go:248`, `mknote/canontags/tags.go:135`, `models/fields.go:338` — unused vars
- `exif/exif.go:733,744,749` — capitalized error strings (ST1005)
- `tiff/tiff.go:52` — `errors.New(fmt.Sprintf(...))` should be `fmt.Errorf(...)` (S1028)
- `tiff/tiff.go:58` — Yoda condition (ST1017)
- `xmp/xmp_test.go:4` — `io/ioutil` deprecated since Go 1.19 (SA1019)

## Pre-existing TODOs already in the source

Not re-described here, just pointers: `exif/exif.go:371,382` (timezone parsing gaps —
GPS time, Nikon WorldTime), `mknote/nikon.go:30` ("fix regression test"),
`tiff/tag.go:178` (tag values pointing outside the current IFD block aren't handled —
common for thumbnails/large blobs per the comment there).

## exif_test.go's TestDecode assumes every fixture has EXIF

`exif/exif_test.go`'s `TestDecode` iterates every `.jpg` in the shared testdata directory
and fails if `Decode` doesn't return usable EXIF. It currently fails on
`imgphash_cat_sky.jpg`, `imgphash_cat_medium.jpg`, and `imgphash_cat_smiling.jpg` — real
photos with no EXIF at all, added to the shared `dedup-testdata` corpus when `dedup`'s
`testdata/` and `large_testdata/` were merged. Pre-existing fallout from that merge, not
from anything in this repo; needs either an explicit skip for known-no-EXIF fixtures or a
different assertion.

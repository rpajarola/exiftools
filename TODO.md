# Known issues

Found via a correctness-focused review of the core parsing packages (`tiff/`, `exif/`,
`mknote/`), plus `go vet`/`staticcheck`. All four bugs below are verified against the
actual source, not just theorized.

## Crash / DoS bugs in the parsing path

These trigger on crafted or corrupted input, not just normal camera files — worth fixing
before parsing untrusted uploads.

### 1. Integer overflow in tag-size bounds check allows multi-GB allocation from a ~12-byte tag

`tiff/tag.go:161-165`: `valLen := size * t.Count` multiplies two `uint32`s with no
overflow check, then compares `valLen` against `TagLengthCutoff` (4MB) to reject
oversized tags. For a multi-byte type (e.g. `DTRational`, size 8) with
`Count = 0x20000000` (536,870,912), `valLen` wraps to `0`, sailing past the cutoff.
`convertVals()` (`tiff/tag.go:359`, `373`, and similarly `300`/`310`/`320`/`330`/`340`/
`350`/`387`/`397` for the other multi-value types) then allocates
`make([][]int64, int(t.Count))` using the raw, unwrapped `Count` — ~536M slice headers,
~12.9GB. Runs unconditionally on every `tiff.Decode()` call. Affects any type with
size ≥ 2 (Short/Long/Float/SLong/Rational/SRational/Double); the size-1 types
(Byte/Ascii/SByte/Undefined) can't wrap this way.

Fix: do the overflow check before multiplying (e.g. `t.Count > TagLengthCutoff/size`),
or use `uint64` arithmetic for the comparison.

### 2. IFD cycle detection only catches an immediate back-reference, not a longer cycle

`tiff/tiff.go:71-95`: the loop tracks only `prev` (the previous IFD offset) and aborts
if `offset == prev`. A 2-cycle (IFD A → IFD B → IFD A → IFD B → ...) never repeats "the
immediately previous" offset, so the check never fires — `t.Dirs = append(t.Dirs, d)`
runs forever, growing unbounded until OOM. `fingerprint/testdata/corrupt/infinite_loop_exif.jpg`
in the consuming `dedup` repo suggests a cycle case was tested at some point; worth
re-verifying it's actually a 2+ cycle and not just a self-loop, since a self-loop is the
one case this check does catch.

Fix: track a set of seen offsets, not just the last one.

## Unrecovered panics on malformed-but-plausible tag data

Both reachable through normal public API calls, not deep internals.

### 3. `mknote/canon.go:134` (`processCameraSettingsMap`) indexes a static table with a raw file value

`return CanonCameraSettingsFields[i][a]` — `a` comes straight from `tag.Int(i)`, i.e.
whatever value the file's CanonCameraSettings tag declares, with no bounds check against
`len(CanonCameraSettingsFields[i])`. A legitimate-looking but out-of-catalog value (e.g.
`ContinuousDrive = 200`) panics with an unrecovered slice-index-out-of-range. Reachable
via `CanonRaw.Get()`, used unconditionally whenever a caller requests Canon structured
data.

### 4. `exif/exif.go:458` (`parse3Rat2`) calls `tag.Rat2(i)` before checking `tag.Count`

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

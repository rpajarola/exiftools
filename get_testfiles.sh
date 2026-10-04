#!/bin/bash
# Fetches testdata/ from the shared rpajarola/dedup-testdata release, which
# also backs rpajarola/dedup's fingerprint tests.
set -euo pipefail

mkdir -p "$(dirname "$0")/testdata"
cd "$(dirname "$0")/testdata"
curl -L -o testdata-images.tar.gz \
  https://github.com/rpajarola/dedup-testdata/releases/latest/download/testdata-images.tar.gz
tar -xzf testdata-images.tar.gz
rm testdata-images.tar.gz

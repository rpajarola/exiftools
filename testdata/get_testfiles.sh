#!/bin/bash

git annex init || echo "git annex not initialized"

urlbase=https://cave.servium.ch/github/exiftools
git annex list | while read where fname; do
  case "${where}:${fname}" in
  *X:testdata/*)
    echo "${fname}"
    git annex addurl --file "${fname}" "${urlbase}/${fname}"
    ;;
  esac
done


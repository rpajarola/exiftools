#!/bin/bash

git annex init || echo "git annex not initialized"
git annex pull

urlbase=https://cave.servium.ch/github/exiftools
git annex list | while read where fname; do
  case "${where}:${fname}" in
  ??_??:*)
    git annex addurl --file "${fname}" "${urlbase}/${fname}"
  esac
done


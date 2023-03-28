#!/bin/bash
#
# Ref. https://superuser.com/a/1012535

# PNG >= 512x512
SRC=$1
DST=$2

tmpd="$(mktemp -d)"
cp "${SRC}" "${tmpd}/icon.png"
pushd "${tmpd}" || exit 0
	convert -scale 16  icon.png icon-16.png
	convert -scale 24  icon.png icon-24.png
	convert -scale 32  icon.png icon-32.png
	convert -scale 48  icon.png icon-48.png
	convert -scale 64  icon.png icon-64.png
	convert -scale 96  icon.png icon-96.png
	convert -scale 128 icon.png icon-128.png
	convert -scale 256 icon.png icon-256.png
	convert -scale 512 icon.png icon-512.png
	convert icon-16.png icon-24.png icon-32.png icon-48.png icon-64.png icon-96.png icon-128.png icon-256.png icon-512.png icon.ico
popd || exit 0
mv "${tmpd}/icon.ico" "${DST}"
rm -fr "${tmpd}"

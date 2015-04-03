#! /bin/bash
#
# d_build_cross.sh
# Copyright (C) 2015 oliveagle <oliveagle@gmail.com>
#
# Distributed under terms of the MIT license.
#

for GOOS in darwin linux windows; do
  for GOARCH in 386 amd64; do
    echo "building $GOOS $GOARCH"
    GOOS=$GOOS GOARCH=$GOARCH go build -v -o bin/hickwall-$GOOS-$GOARCH
    echo ""
  done
done

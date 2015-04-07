#! /bin/bash
#
# d_build_cross.sh
# Copyright (C) 2015 oliveagle <oliveagle@gmail.com>
#
# Distributed under terms of the MIT license.
#

for GOOS in darwin linux; do
  for GOARCH in 386 amd64; do
    echo "building $GOOS $GOARCH"
    GOOS=$GOOS GOARCH=$GOARCH go build -v -o bin/hickwall-$GOOS-$GOARCH
    echo ""
  done
done

# windows only support 386 version. don't know way amd64 cannot be compiled.
GOOS=windows GOARCH=386 go build -v -o bin/hickwall-windows-386.exe

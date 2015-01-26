#! /bin/bash
#
# build.sh
# Copyright (C) 2015 oliveagle <oliveagle@gmail.com>
#
# Distributed under terms of the MIT license.
#


mkdir -p bin
echo "building windows 386 into bin/  ..." && \
  GOOS=windows GOARCH=386 go build -o bin/hickwall.exe && \
echo "building linux amd64 into bin/  ..." && \
  GOOS=linux GOARCH=amd64 go build -o bin/hickwall.linux.amd64 && \
echo "building darwin amd64 into bin/ ..." && \
  GOOS=darwin GOARCH=amd64 go build -o bin/hickwall.darwin.amd64


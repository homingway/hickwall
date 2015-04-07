#! /bin/bash
#
# build_all.sh
# Copyright (C) 2015 oliveagle <oliveagle@gmail.com>
#
# Distributed under terms of the MIT license.
#
SCRIPT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

APP_ROOT="$SCRIPT_ROOT/.."
APP_NAME="hickwall"
GOOS="windows"
GOARCH=386
BUILD_CMD="go build -v -o bin/hickwall-$GOOS-$GOARCH.exe"
GOIMG="golang:1.4.2-cross"
GOPATH="/oledev/gocodez/"

docker run --rm \
  -v $APP_ROOT:/usr/src/$APP_NAME -w /usr/src/$APP_NAME \
  -v $GOPATH:/gopath -e GOPATH=/gopath \
  -e GOOS=$GOOS -e GOARCH=$GOARCH \
  $GOIMG $BUILD_CMD

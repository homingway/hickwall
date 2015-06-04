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

VER=$(cat $APP_ROOT/release-version)

cd $SCRIPT_ROOT

$(bash make_win_helper.sh)

cd $APP_ROOT

GIT_HASH=$(git rev-parse --short HEAD)
# save this hash for later packing
echo "$GIT_HASH" > "$SCRIPT_ROOT/GIT_HASH"

echo "Version: $VER, Build: $GIT_HASH"

go build -ldflags "-X main.Version $VER -X main.Build $GIT_HASH" -v -o hickwall.exe && cp hickwall.exe bin/hickwall-windows-386.exe

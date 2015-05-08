#! /bin/bash
#
# build_all.sh
# Copyright (C) 2015 oliveagle <oliveagle@gmail.com>
#
# Distributed under terms of the MIT license.
#
SCRIPT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

APP_ROOT="$SCRIPT_ROOT/.."
VER=$(cat $APP_ROOT/release-version)

cd $APP_ROOT

sed -i "s/\(app.Version.*\"\)\(.*\)\"$/\1v$VER\"/g" hickwall.go

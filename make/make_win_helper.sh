#! /bin/bash
#
# build_all.sh
# Copyright (C) 2015 oliveagle <oliveagle@gmail.com>
#
# Distributed under terms of the MIT license.
#
SCRIPT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

APP_ROOT="$SCRIPT_ROOT/.."

cd "$APP_ROOT/misc/win_helper_service"
gcc -Os hickwall_helper.c -o hickwall_helper.exe && mv hickwall_helper.exe "$APP_ROOT/bin/"
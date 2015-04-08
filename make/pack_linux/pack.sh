#! /bin/bash
#
# pack_with_fpm.sh
# Copyright (C) 2015 oliveagle <oliveagle@gmail.com>
#
# Distributed under terms of the MIT license.
#


SCRIPT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJ_ROOT="$SCRIPT_ROOT/../.."
mkdir -p "$PROJ_ROOT/bin/dist"

docker run --rm -ti -v /oledev:/oledev -w /oledev/gocodez/src/github.com/oliveagle/hickwall lonefreak/fpm bash make/pack_linux/pack_with_fpm.sh


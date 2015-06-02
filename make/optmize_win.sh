#! /bin/bash
#
# build_win.sh
# Copyright (C) 2015 oliveagle <oliveagle@gmail.com>
#
# Distributed under terms of the MIT license.
#
SCRIPT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

cd "$SCRIPT_ROOT/.." && rm -f hickwall.exe  && go build -gcflags -m 2>&1 | tee gcflags_m.log


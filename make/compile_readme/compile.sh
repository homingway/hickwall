#! /bin/bash
#
# compile.sh
# Copyright (C) 2015 oliveagle <oliveagle@gmail.com>
#
# Distributed under terms of the MIT license.
#


docker run -v /oledev:/oledev -w /oledev/gocodez/src/github.com/oliveagle/hickwall python_grip grip --wide --export Readme.md


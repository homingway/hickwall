#! /bin/bash
#
# build.sh
# Copyright (C) 2015 oliveagle <oliveagle@gmail.com>
#
# Distributed under terms of the MIT license.
#


docker build -t wine_innosetup .

#docker run --rm -ti -v /oledev:/oledev -w /oledev/gocodez/src/github.com/oliveagle/hickwall  wine_1 bash make/pack.win


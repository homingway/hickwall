#! /bin/bash
#
# pack_wi_inno.sh
# Copyright (C) 2015 oliveagle <oliveagle@gmail.com>
#
# Distributed under terms of the MIT license.
#


docker run --rm -ti -v /oledev:/oledev -w /oledev/gocodez/src/github.com/oliveagle/hickwall wine_innosetup bash make/pack_win/pack_with_inno.sh


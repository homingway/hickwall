
#
# after_install.sh
# Copyright (C) 2015 oliveagle <oliveagle@gmail.com>
#
# Distributed under terms of the MIT license.
#

VER=

if [ -z "$VER" ]; then
  echo "VER is empty" 1>&2
  exit 1
fi

rm -f /opt/hickwall/current
ln -s /opt/hickwall/$VER /opt/hickwall/current

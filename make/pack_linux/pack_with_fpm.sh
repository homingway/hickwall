#! /bin/bash
#
# pack_with_fpm.sh
# Copyright (C) 2015 oliveagle <oliveagle@gmail.com>
#
# Distributed under terms of the MIT license.
#

SCRIPT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJ_ROOT="$SCRIPT_ROOT/../.."
DIST="$PROJ_ROOT/bin/dist"
VER=$(cat $PROJ_ROOT/release-version)

mkdir -p /opt/hickwall/$VER
mkdir -p /opt/hickwall/shared

cp $PROJ_ROOT/config.yml /opt/hickwall/shared/config.yml.example
cp $PROJ_ROOT/Readme.html /opt/hickwall/$VER/
cp $PROJ_ROOT/Readme.md /opt/hickwall/$VER/


# amd64
cp -f $PROJ_ROOT/bin/hickwall-linux-amd64 /opt/hickwall/$VER/hickwall
fpm -s dir -t rpm -p $DIST -a x86_64 \
  --version $VER --force --name "hickwall" \
  /opt/hickwall/shared/config.yml.example /opt/hickwall/$VER/

fpm -s dir -t deb -p $DIST -a x86_64 \
  --version $VER --force --name "hickwall" \
  /opt/hickwall/shared/config.yml.example /opt/hickwall/$VER/

# 386
cp -f $PROJ_ROOT/bin/hickwall-linux-386 /opt/hickwall/$VER/hickwall
fpm -s dir -t rpm -p $DIST -a x86 \
  --version $VER --force --name "hickwall" \
  /opt/hickwall/shared/config.yml.example /opt/hickwall/$VER/

fpm -s dir -t deb -p $DIST -a x86 \
  --version $VER --force --name "hickwall" \
  /opt/hickwall/shared/config.yml.example /opt/hickwall/$VER/

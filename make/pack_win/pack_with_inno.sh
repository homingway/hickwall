#!/bin/sh

# get current running script location
SCRIPT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJ_ROOT="$SCRIPT_ROOT/../.."
ISCC=$HOME/bin/iscc
VER=$(cat $PROJ_ROOT/release-version)
DST="$PROJ_ROOT/bin/dist/hickwall-setup-$VER.exe"

TMP_DIR=$(mktemp -d)
echo "temp dir: " $TMP_DIR

cd "$SCRIPT_ROOT"
cp win.iss $TMP_DIR/
cp start.bat $TMP_DIR/
cp stop.bat $TMP_DIR/

cd $TMP_DIR
sed -i "s/{{version}}/$VER/g" win.iss
sed -i "s/{{version}}/$VER/g" start.bat
sed -i "s/{{version}}/$VER/g" stop.bat
# cat win.iss

cd "$PROJ_ROOT"
cp bin/hickwall-windows-386.exe $TMP_DIR/hickwall.exe && \
cp bin/hickwall_helper.exe $TMP_DIR/hickwall_helper.exe && \
  cp config.yml $TMP_DIR/config.yml.example && \
  cp core_config.yml $TMP_DIR/core_config.yml.example && \
  cp Readme.md $TMP_DIR/ && \
  cp Readme.html $TMP_DIR/ && \
  cd $TMP_DIR && \
  $ISCC win.iss && \
  cp Output/setup.exe $DST && \
  echo "copied setup into $DST"


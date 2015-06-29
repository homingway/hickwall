#! /bin/bash
#
# test_all.sh
# Copyright (C) 2015 oliveagle <oliveagle@gmail.com>
#
# Distributed under terms of the MIT license.
#
# get current running script location
SCRIPT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
VENDER_ROOT="$SCRIPT_ROOT/_vendor"

if [ $OS == "Windows_NT" ];then
  export GOPATH=$(cygpath -w "$VENDER_ROOT")
else
  export GOPATH=$VENDER_ROOT
fi

echo "GOPATH=$GOPATH"

#go test ./...
go test ./... -v | grep -E "(--- FAIL)|(^FAIL\s+)|(^ok\s+)"

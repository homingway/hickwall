#! /bin/bash
#
# test_all.sh
# Copyright (C) 2015 oliveagle <oliveagle@gmail.com>
#
# Distributed under terms of the MIT license.
#
# get current running script location
# SCRIPT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
# VENDER_ROOT="$SCRIPT_ROOT/_vendor"

# export GOPATH=$VENDER_ROOT
# if [ $OS == "Windows_NT" ] && [ -n "$(which cygpath)" ];then
#         export GOPATH=$(cygpath -w "$VENDER_ROOT")
# fi

# echo "GOPATH=$GOPATH"

go test ./... -v
# go test ./... -v | grep -E "(--- FAIL)|(^FAIL\s+)|(^ok\s+)"

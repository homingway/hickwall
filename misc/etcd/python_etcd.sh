#! /bin/bash
#
# python_etcd.sh
# Copyright (C) 2015 oliveagle <oliveagle@gmail.com>
#
# Distributed under terms of the MIT license.
#


docker run --rm -ti -v $(pwd):/src  -w /src ole-python-etcd python $@


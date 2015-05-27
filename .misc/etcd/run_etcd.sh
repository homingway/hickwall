#! /bin/bash
#
# run_etcd.sh
# Copyright (C) 2015 oliveagle <oliveagle@gmail.com>
#
# Distributed under terms of the MIT license.
#

docker run -d -p 4001:4001 -p 7001:7001 -v /oledev/var/data/etcd_data:/data 192.168.81.136:5000/microbox/etcd -name etcd-1 


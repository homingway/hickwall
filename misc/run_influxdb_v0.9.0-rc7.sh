#! /bin/bash
#
# run.sh
# Copyright (C) 2015 oliveagle <oliveagle@gmail.com>
#
# Distributed under terms of the MIT license.
#
docker run -d -p 80:80 -p 8083:8083 -p 8084:8084 -p 8086:8086 -p 9022:22 --name="influxdb_090" influxdb:v0.9.0-rc7 /usr/bin/supervisord


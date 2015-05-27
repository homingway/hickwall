#! /usr/bin/env python
# -*- coding: utf-8 -*-
# vim:fenc=utf-8
#
# Copyright © 2015 oliveagle <oliveagle@gmail.com>
#
# Distributed under terms of the MIT license.
#
# File:  post_data.py
# Date:  2015-04-28 07:42

"""

"""

import yaml
import json
import etcd

PATH = "/config/host/DST54869.json"

yml_raw = open("config.yml", "r").read()
ydata = yaml.load(yml_raw)

client = etcd.Client(host="10.0.2.15", port=4001)

y_dumps = json.dumps(ydata)

client.write(PATH, y_dumps)

dumps = client.get(PATH).value
jdata = json.loads(dumps)

print jdata

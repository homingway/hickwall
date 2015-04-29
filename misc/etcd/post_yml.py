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

PATH = "/config/host/DST54869.yml"

yml_raw = open("config.yml", "r").read()

client = etcd.Client(host="10.0.2.15", port=4001)

client.write(PATH, yml_raw)

dumps = client.get(PATH).value
print yaml.load(dumps)

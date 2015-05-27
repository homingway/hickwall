import yaml
import json
import etcd

PATH = "/config/host/DST54869.yml"

yml_raw = open("config.yml", "r").read()

client = etcd.Client(host="10.0.2.15", port=4001)

dumps = client.get(PATH).value
print yaml.load(dumps)

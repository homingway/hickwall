Hickwall
==========

A metric collection and reporting daemon for major platforms. 

***Under heavy construction!***


##Build From Source

You need install [`tools/godep`][url_godep] first. which is a golang dependencies management system. It saved all third-party dependencies under `Godeps`. 

	# build project from source
	godep go build .
	
	# cross build windows binary
	GOOS=windows GOARCH=amd64 godep go build .


##Usage


	# print help info
	hickwall help

	sudo hickwall install

	sudo hickwall start


##Configuration



## Development

currently we support both influxdb v0.9.0-rc7 and v0.8.8. while developing. you have to copy and paste `"github.com/influxdb/influxdb"` to `"github.com/influxdb/influxdb_088"` and then `checkout -b v0.8.8` in `influxdb_088` folder.




[url_godep]: https://github.com/tools/godep "tools/godep"
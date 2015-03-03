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



[url_godep]: https://github.com/tools/godep "tools/godep"
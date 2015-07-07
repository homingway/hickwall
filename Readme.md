Hickwall
==========
A metric collecting and reporting daemon for major platforms. Currently we are focusing on windows.

[![Build status](https://ci.appveyor.com/api/projects/status/o8mfdexkpebe60k6?svg=true)](https://ci.appveyor.com/project/oliveagle/hickwall)
we are mocking services that we are still relay on in unittests.


***Under heavy construction!***


##Build From Source

You need install [`mattn/gom`][url_gom] first. which is a golang dependencies management system. It saved all third-party dependencies under `_vendor`. You can also download any dependencies and put it into "_vendor" yourself.

    go get github.com/mattn/gom
    gom install
    # copy influxdb into influxdb_088 and checkout v0.8.8
    gom test ./... -v
    gom build

##Usage

	# print help info
	hickwall help

	sudo hickwall install

	sudo hickwall start


##Configuration

there are **three ways** to config hickwall client.
* Use local configuration file:  hickwall can run standalone with all configurations in `shared/config.yml`

* Use Remote configuration service: hickwall can also retrive configuration from `etcd` cluster without encryption. but you have to write a minimal `shared/config.yml` to tell hickwall where to find etcd.

    ```yaml
    config_strategy: "etcd"
    etcd_machines:
        - "http://127.0.0.1:4001"
    etcd_path: "/config/host/myhost.yml"
    ```
* Use Contral Registry to config. we are working on an contral registration service.

	```yaml
    config_strategy: "registry"
    registry_urls:
        - "http://127.0.0.1:8080/agent_registry"
    # api listen port, default is 3031
    listen_port: 3031
    ```

## Development

currently we support both influxdb v0.9.0 and v0.8.8. while developing. you have to copy and paste `"github.com/influxdb/influxdb"` to `"github.com/influxdb/influxdb_088"` and then `checkout -b v0.8.8` in `influxdb_088` folder.




[url_gom]: https://github.com/mattn/gom "mattn/gom"

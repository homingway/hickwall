### Version 0.1.0
* windows service install/remove/start/stop/restart
* easy development environment setup using docker. 
    * compile readme.md to readme.html using python grip in docker
    * cross compile windows/linux binaries in docker
    * pack rpm/deb/windows setup in docker
* serveral backends supported:
    * influxdb v0.8.8
    * influxdb v0.9.0-rc7
* collectors implemented:
    * windows_pdh
    * windows_wmi

### Version 0.1.1
* config devided into 2 part. 1st one is core config, 2nd is runtime config.
* runtime config can be loaded from file or from etcd
* runtime config can be hot reloaded by watching etcd path changes.
* collector_cmd implemented
* collector_ping implemented
* refactoried how collectors been created.
* fixed several bugs
* added heartbeat metrics
* removed hickwallhelper service on windows.

### Version 0.2.1
* collector_cmd removed for risk consideration
* use c implemented hickwallhelper service on windows.
* implemented a newcore. 
* dummy backend
* file backend
* kafka backend
* kafka subscriber. (e.g. 1(or few) kafka subscriber, N kafka producer. )
* remove 3rd-party logging module. replaced with a customized wrapper of builtin log
* add RSS limit of the hickwall process.
* unittest covered most packages.

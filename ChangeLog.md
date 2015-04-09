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

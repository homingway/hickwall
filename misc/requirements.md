
####client端的需求：
* 服务控制：install, uninstall, start, stop, restart, (auto start on boot)
* 服务配置：
  * 全局tags
  * 日志记录的级别，格式，rotate size，rotate 个数等
  * 客户端自己的metrics是否上报，上报频率
  * 客户端heartbeat
  * transport配置：
    * queue配置， backfill支持。
    * 支持的transport的配置： influxdb, graphite, amqp...
  * collecotr配置：
    * 默认时间间隔
    * win_pdh 配置：
      * interval
      * tags
      * queries
    * win_wmi 配置：
      * interval
      * tags
      * queries
    * 支持自定义脚本执行采集：定义脚本返回的格式。


####配置管理服务

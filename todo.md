# TODO:
---------------
* collector host basic info and send back with heartbeat




# Postponed: 
---------------

* collector aggregation
    
```` yaml
    # TODO: [postponed 2015Mar05] collector_aggregation would be better to have.
    # `Exits` aggregation function can apply on any kind of value.  
    # But in most cases, only `int` and `float` can be aggregated. Internally, all numbers
    # will be converted into float64
    # TODO: graphite aggregator like aggregation ??
    collector_aggregation:
        - 
            # TickDriven
            # single source and simple calc
            Interval: 60s
            source: {
                "A": "win.phd.memory.available_bytes"
            }
            aggregation: "A / 1024 / 1024"
            metric: "win.pdh.memory.available.mb"
            inherit_tags: "A"
            tags: {
                "new_tag": "aggregation"
            }
            # Merge(A.Tags).Merge(tags)
            inherit_meta: "A"
            meta: {
                "unit": "MegaBytes",
                "interval": "60s"
            }

        - 
            # aggregated on tags
            Interval: 60s
            source: {
                "A": "win.pdh.memory.available.mb"
            }

            # because we are aggregated on tags. so tags been aggregated will not been merged
            source_tags: 
                - "bu"

            aggregation: "Sum(A)"
            metric: "win.pdh.memory.available.mb.{{.bu}}"
            inherit_meta: "A"

        - 
            # TickDriven
            # source can be generated from other aggregations
            # single source and predefined aggregation functions with tick period
            source: {
                "A": "win.pdh.memory.available.mb"
            }
            aggregation: "Average(WindowTick(A, 5))"
            metric: "win.pdh.memory.available.mb.avg.5tick"

        - 
            # single source with '*'
            source: {
                "A": "win.wmi.fs.*.size.bytes"
            }

            # for `Sum` aggregation function

            # s1     s2     s3     s4     s5
            # --------------------------------
            # A1.1   A1.2   A1.3   A1.4   A1.5
            # -      A2.1   -      A2.2   - 
            # -      -      A3.1   -      -

            # s1 = A1.1
            # s2 = A1.2 + A2.1
            # s3 = A1.3 + A2.1 + A3.1

            # S = Last(A1) + Last(A2) + Last(A3)

            aggregation: "Sum(A) / 1024 / 1024"
            metric: "win.wmi.fs.total_size.mb"
            # if drop_source == true, source metric data will not be sent to transport middleware.
            drop_source: true

        - 
            # TODO: multiple source value may not arrived at the same time. 

            # multiple source and simple calc
            source: {
                "A": "win.phd.memory.available_bytes",
                "B": "win.wmi.mem.totalphysicalmemory"
            }

            # for `Percent` aggregation function, value will only be avaible if the last 
            # data source have valid data tick comes in.
            # `A` is srouce A, `B` is source B, `P` is Percent

            # A A A A A A A 
            # - - - - B - - 
            # - - - - P P P

            # TODO: [done 2015mar05]how long we have to keep data `B`, forever in stats file

            # TODO: persist program states forever.
            # Last Data point B will be persist forever.
            aggregation: "Percent(A, Last(B))"  

            metric: "win.pdh.memory.available.pct"
            # Merge(B.Tags).Merge(A.Tags)
            inherit_tags: "B,A" 

        - 
            # single source and perdefined aggregation functions with minute period

            # TODO: aggregation on Minute Time Window
            # Note: aggregation on Minute time window is a little bit tricky. data
            # will not be generated once a new data point comes in. coz we may not
            # know if the new data point is the last data tick within the unit(Minute, Day, ...)
            #  of this time window. So we have to wait the time unit to close. then 
            # we can calculate and emit the data.

            # `Average(WindowMinute(A, 3))` calculation process is shown as below.
            # m1 ... m5 represent five minutes. `[- -]` represent the minute has
            # two time span, the first half minute, and the second half minute.
            # `@` represent to a data tick comes in.  `|` and beneath represent
            # a aggregation happened at that time.

            #  m1    m2    m3    m4    m5
            # [- -] [- -] [- -] [- -] [- -]
            #  @     @       @
            #                    |
            #                    Average(m1, m2, m3)
            #                    @ 
            #                          |
            #                          Average(m2, m3, m4)     
            #                            @
            # [- -] [- -] [- -] [- -] [- -]

            source: {
                "A": "win.pdh.memory.available.mb"
            }
            aggregation: "Average(WindowMinute(A, 5))"
            metric: "win.pdh.memory.available.mb.avg.5min"

        -
            # single source, exists
            source: {
                "A": "win.wmi.service.iis"
            }
            Interval: 60s
            aggregation: "Exists(Last(A))"
            aggregation: "Exists(Last(WindowTick(A, 100)))"
            aggregation: "Exists(Last(WindowHour(A, 24)))"
            metric: "win.wmi.service.iis.installed"

        - 
            # TickDriven
            # single source and delta
            source: {
                "A": "win.phd.app.visits"
            }
            # delta can only applied on Window.size == 2
            aggregation: "Delta(WindowTick(A,2))"
            metric: "win.pdh.app.visits.delta.tick"
        - 
            # TickDriven
            # single source and delta
            source: {
                "A": "win.phd.app.visits"
            }
            # delta can only applied on Window.size == 2
            aggregation: "Delta(WindowMinute(A,2))"
            metric: "win.pdh.app.visits.delta.minute"

        - 
            # single source, delta
            source: {
                "A": "win.pdh.app.visits.delta.minute"
            }
            aggregation: "Average(WindowMinute(A, 5))"
            metric: "win.pdh.app.visits.delta.avg.5m"

            # if drop_source == true, source metric data will not be sent to transport middleware
            drop_source: true
````

* collector_mysql
````yaml
# TODO: collector_mysql_query
collector_mysql_query:
    - 
        tags: [
            [ "bu", "test" ],
        ]
        host     : "127.0.0.1"
        port     : 3306
        username : "root"
        password : "root"

        queries  :
            -
                metric_key: "mysql_query.xxxxx"
                # tags: [
                #     [ "some", "test2"]
                # ]

                database: "db1"
                desc: "one line query"
                query: "SELECT count(*) as cnt, sum(*) as total FROM sometable where column=123"
                values_from: "cnt"

            # metrics:
            #   mysql_query.xxxx  {"bu":"test", "some", "test2"}      123

            -
                database: "db1"
                desc: "multiple line query"
                # multiple line is tricky. indent is very important
                #
                #  xxxx: >
                #      a          =>  "a b\n"
                #      b
                #  xxxx: >
                #      a          =>  "a\n b\n"
                #       b
                query: >
                    SELECT count(*) as cnt, sum(*) as
                    total FROM sometable where column=123

                metric_key: "mysql_query.xxxxx"

                # multiple values from also ok
                values_from: "cnt, total"

            # metrics:
            #   mysql_query.xxxx.cnt    {"bu": "test"}  123
            #   mysql_query.xxxx.total  {"bu": "test"}  432
````

* use wal and replay to replace boltq
* transport_amqp
````yaml
# TODO: support amqp transport
# transport_amqp_enabled: false
# transport_amqp_hosts: ["127.0.0.1:5672", "127.0.0.1:5672"]
# transport_amqp_vhost: "/"
# transport_amqp_username: "guest"
# transport_amqp_password: "guest"
# transport_amqp_exchange: "amq.fanout"
# transport_amqp_exchange_type: "fanout"
# transport_amqp_routing_key: "hickwall"
# transport_amqp_persistent: false
# transport_amqp_queue: "queuename"
````


# DONE
* collector_ping (20150507)
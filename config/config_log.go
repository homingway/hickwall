package config

import (
	"fmt"
	"github.com/oliveagle/hickwall/_third_party/stringio"
	l4g "github.com/oliveagle/log4go"
	"github.com/spf13/viper"
	"text/template"
)

func parseLoggerConfig() {

	viper.SetDefault("Consolelog.Level", "DEBUG")
	Conf.setDefault("Consolelog.Level", "DEBUG")
	Conf.setDefault("Consolelog.Format", "[%D %T] [%L] (%S) %M")

	Conf.setDefault("Filelog.Level", "DEBUG")

	// TODO: windows default logfile
	Conf.setDefault("Filelog.Logfile", "/var/log/hickwall/hickwall.log")
	Conf.setDefault("Filelog.Format", "[%D %T] [%L] (%S) %M")
	Conf.setDefault("Filelog.Rotate", false)
	Conf.setDefault("Filelog.Maxsize", "0M")
	Conf.setDefault("Filelog.Maxlines", "0K")
	Conf.setDefault("Filelog.Daily", false)

	l4g.LoadConfigurationFromString(`<logging>
	  <filter enabled="true">
	    <tag>stdout</tag>
	    <type>console</type>
	    <!-- level is (:?FINEST|FINE|DEBUG|TRACE|INFO|WARNING|ERROR) -->
	    <level>DEBUG</level>
	  </filter>
	  </logging>`)

	// tpl := `<logging>
	//      <filter enabled="true">
	//        <tag>stdout</tag>
	//        <type>console</type>
	//        <!-- level is (:?FINEST|FINE|DEBUG|TRACE|INFO|WARNING|ERROR) -->
	//        <level>{{.Consolelog.level}}</level>
	//        <format>{{.Consolelog.format}}</format>
	//      </filter>
	//        <filter enabled="true">
	//    <tag>file</tag>
	//    <type>file</type>
	//    <level>{{.Level}}</level>
	//    <property name="filename">test.log</property>
	//    <!--
	//       %T - Time (15:04:05 MST)
	//       %t - Time (15:04)
	//       %D - Date (2006/01/02)
	//       %d - Date (01/02/06)
	//       %L - Level (FNST, FINE, DEBG, TRAC, WARN, EROR, CRIT)
	//       %S - Source
	//       %M - Message
	//       It ignores unknown format strings (and removes them)
	//       Recommended: "[%D %T] [%L] (%S) %M"
	//    -->
	//    <!--property name="format">[%D %T] [%L] (%S) %M</property-->
	//    <property name="format">{{.Format}}</property>
	//    <property name="rotate">{{.Rotate}}</property> <!-- true enables log rotation, otherwise append -->
	//    <property name="maxsize">{{.Maxsize}}</property> <!-- \d+[KMG]? Suffixes are in terms of 2**10 -->
	//    <property name="maxlines">{{.Maxlines}}</property> <!-- \d+[KMG]? Suffixes are in terms of thousands -->
	//    <property name="daily">{{.Daily}}</property> <!-- Automatically rotates when a log message is written after midnight -->
	//  </filter>
	//      </logging>
	//    `
	tpl := `<logging>
          <filter enabled="true">
            <tag>stdout</tag>
            <type>console</type>
            <!-- level is (:?FINEST|FINE|DEBUG|TRACE|INFO|WARNING|ERROR) -->
            <level>{{.Consolelog.Level}}</level>
            <format>{{.Consolelog.Format}}</format>
          </filter>
           <filter enabled="true">
           <tag>file</tag>
           <type>file</type>
           <level>{{.Level}}</level>
           <property name="filename">test.log</property>
           <!--
              %T - Time (15:04:05 MST)
              %t - Time (15:04)
              %D - Date (2006/01/02)
              %d - Date (01/02/06)
              %L - Level (FNST, FINE, DEBG, TRAC, WARN, EROR, CRIT)
              %S - Source
              %M - Message
              It ignores unknown format strings (and removes them)
              Recommended: "[%D %T] [%L] (%S) %M"
           -->
           <!--property name="format">[%D %T] [%L] (%S) %M</property-->
           <property name="format">{{.Format}}</property>
           <property name="rotate">{{.Rotate}}</property> <!-- true enables log rotation, otherwise append -->
           <property name="maxsize">{{.Maxsize}}</property> <!-- \d+[KMG]? Suffixes are in terms of 2**10 -->
           <property name="maxlines">{{.Maxlines}}</property> <!-- \d+[KMG]? Suffixes are in terms of thousands -->
           <property name="daily">{{.Daily}}</property> <!-- Automatically rotates when a log message is written after midnight -->
        </filter>
    </logging>`
	// fmt.Println(tpl)

	sio := stringio.NewStringIO()
	t := template.Must(template.New("logconf").Parse(tpl))
	err := t.Execute(sio, Conf)
	if err != nil {
		fmt.Printf("err: %v", err)
	}
	fmt.Println(sio.GetValueString())
	l4g.Info("hahahah")
}

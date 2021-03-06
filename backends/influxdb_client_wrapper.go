package backends

import (
	"fmt"
	client090 "github.com/influxdb/influxdb/client"
	"github.com/influxdb/influxdb/influxql"
	client088 "github.com/influxdb/influxdb_088/client"
	"github.com/oliveagle/hickwall/logging"
	"github.com/oliveagle/hickwall/newcore"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	pat_influxdb_version = regexp.MustCompile(`[v]?\d+\.\d+\.\d+[\S]*`)
)

type InfluxdbClient interface {
	Version() string
	Ping() (time.Duration, string, error)
	Write(client090.BatchPoints) (*client090.Response, error)
	Query(client090.Query) (*client090.Response, error)
	IsCompatibleVersion(v string) bool
}

func influxdbParseVersionFromString(v string) (version string) {
	// parse version and trimleft "v"
	ss := pat_influxdb_version.FindAllString(v, -1)
	if len(ss) > 0 {
		version = ss[0]
		if strings.HasPrefix(version, "v") == true {
			version = strings.TrimLeft(version, "v")
		}
	}
	return
}

func NewInfluxdbClient(conf map[string]interface{}, version string) (InfluxdbClient, error) {
	v := influxdbParseVersionFromString(version)
	if v == "0.9.0" {
		return NewInfluxdbClient_v090(conf)
	} else if v == "0.8.8" {
		return NewInfluxdbClient_v088(conf)
	}
	return nil, fmt.Errorf("incompatible version of influxdb: %s", v)
}

// --------------------------------  version: v0.9.0 --------------------------------------------

type InfluxdbClient_v090 struct {
	client *client090.Client
}

func (c *InfluxdbClient_v090) Version() string {
	return "0.9.0"
}

func (c *InfluxdbClient_v090) Ping() (time.Duration, string, error) {
	return c.client.Ping()
}

func (c *InfluxdbClient_v090) Write(bp client090.BatchPoints) (*client090.Response, error) {
	return c.client.Write(bp)
}

func (c *InfluxdbClient_v090) Query(q client090.Query) (*client090.Response, error) {
	return c.client.Query(q)
}

func (c *InfluxdbClient_v090) IsCompatibleVersion(v string) bool {
	if influxdbParseVersionFromString(v) == c.Version() {
		return true
	}
	return false
}

func NewInfluxdbClient_v090(conf map[string]interface{}) (*InfluxdbClient_v090, error) {
	tmp_conf := map[string]interface{}{}
	for key, value := range conf {
		tmp_conf[strings.ToLower(key)] = value
	}
	url_str, ok := tmp_conf["url"]
	if ok != true {
		return nil, fmt.Errorf("version 0.9.0, config missing: URL")
	}
	host_url, err := url.Parse(url_str.(string))
	if err != nil {
		return nil, fmt.Errorf("version 0.9.0 cannot parse url: %s, err: %v", url_str, err)
	}

	username, ok := tmp_conf["username"]
	if ok != true {
		return nil, fmt.Errorf("version 0.9.0, config missing: Username")
	}

	password, ok := tmp_conf["password"]
	if ok != true {
		return nil, fmt.Errorf("version 0.9.0, config missing: Password")
	}

	useragent, ok := tmp_conf["useragent"]
	if ok != true {
		return nil, fmt.Errorf("version 0.9.0, config missing: UserAgent")
	}

	c := client090.Config{
		URL:       *host_url,
		Username:  username.(string),
		Password:  password.(string),
		UserAgent: useragent.(string),
	}

	cli, err := client090.NewClient(c)
	if err != nil {
		return &InfluxdbClient_v090{}, fmt.Errorf("version 0.9.0, cannot create client: %v", err)
	}
	return &InfluxdbClient_v090{
		client: cli,
	}, nil
}

// --------------------------------  version: v0.8.8 --------------------------------------------

type InfluxdbClient_v088 struct {
	client     *client088.Client
	httpClient *http.Client
	schema     string
	host       string
	username   string
	password   string
	flat_tpl   string
}

func (c *InfluxdbClient_v088) Version() string {
	return "0.8.8"
}

func (c *InfluxdbClient_v088) IsCompatibleVersion(v string) bool {
	if influxdbParseVersionFromString(v) == c.Version() {
		return true
	}
	return false
}

func (c *InfluxdbClient_v088) getUrl(path string) string {
	return fmt.Sprintf("%s://%s%s?u=%s&p=%s", c.schema, c.host, path, c.username, c.password)
}

func (c *InfluxdbClient_v088) Ping() (time.Duration, string, error) {

	now := time.Now()

	url := c.getUrl("/ping")
	// fmt.Println(url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, "", err
	}

	req.Header.Set("User-Agent", "hickwall client")

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return 0, "", err
	}
	version := resp.Header.Get("X-Influxdb-Version")
	defer resp.Body.Close()

	return time.Since(now), version, nil
}

func (c *InfluxdbClient_v088) Write(bp client090.BatchPoints) (*client090.Response, error) {
	// logging.Debug("InfluxdbClient_v088.Write")
	// v0.9.0-rc7 [
	//  {
	//      Name: "a",
	//      Timestamp: "1",
	//      Fields: {"f1": "v1", "f2": "v2"},
	//      Precision: "s"
	//  }
	// ]

	// v0.8.8  [
	//   {
	//     "name": "log_lines",
	//     "columns": ["time", "sequence_number", "line"],
	//     "points": [
	//       [1400425947368, 1, "this line is first"],
	//       [1400425947368, 2, "and this is second"]
	//     ]
	//   }
	// ]

	var series []*client088.Series

	for _, p := range bp.Points {
		s := client088.Series{}
		// s.Name = p.Name
		name, err := newcore.FlatMetricKeyAndTags(c.flat_tpl, p.Measurement, p.Tags)
		if err != nil {
			logging.Error("FlatMetricKeyAndTags Failed!", err)
			return nil, err
		}
		s.Name = name

		point := []interface{}{}

		// time, first
		s.Columns = append(s.Columns, "time")
		point = append(point, p.Time.UnixNano()/1000000)

		// then others
		for key, value := range p.Fields {
			s.Columns = append(s.Columns, key)
			point = append(point, value)
		}

		s.Points = append(s.Points, point)

		logging.Tracef("influxdb --> %+v", s)

		series = append(series, &s)
	}

	// pretty.Println(series)

	err := c.client.WriteSeriesWithTimePrecision(series, "ms")
	if err != nil {
		logging.Errorf("InfluxdbClient_v088.Write.WriteSeriesWithTimePrecision Error: %v", err)
	} else {
		logging.Trace("InfluxdbClient_v088.Write Done No Error")
	}

	return nil, err
}

func (c *InfluxdbClient_v088) Query(q client090.Query) (*client090.Response, error) {
	series, err := c.client.Query(q.Command, "ms")
	// fmt.Println(series, err)

	res := client090.Result{
		Series: []influxql.Row{},
		Err:    fmt.Errorf(""),
	}

	// v0.8.8  [
	//   {
	//     "name": "log_lines",
	//     "columns": ["time", "sequence_number", "line"],
	//     "points": [
	//       [1400425947368, 1, "this line is first"],
	//       [1400425947368, 2, "and this is second"]
	//     ]
	//   }
	// ]

	// type Row struct {
	// 	Name    string            `json:"name,omitempty"`
	// 	Tags    map[string]string `json:"tags,omitempty"`
	// 	Columns []string          `json:"columns"`
	// 	Values  [][]interface{}   `json:"values,omitempty"`
	// 	Err     error             `json:"err,omitempty"`
	// }
	// fmt.Println("-------------------------")
outer:
	for _, ss := range series {
		// fmt.Println(ss.Name)
		// fmt.Println(ss.GetColumns())
		// fmt.Println(ss.GetPoints())

		idx_time := -1
	inner:
		for idx, v := range ss.GetColumns() {
			if v == "time" {
				idx_time = idx
				break inner
			}
		}
		points := ss.GetPoints()

		if idx_time >= 0 {
			// convert time back
			for _, point := range points {
				//TODO: here maybe problematic
				ms, err := client090.EpochToTime(int64(point[idx_time].(float64)), "ms")
				if err != nil {
					points = ss.GetPoints()
					break outer
				}
				point[idx_time] = ms
			}
		}

		row := influxql.Row{
			Name:    ss.GetName(),
			Columns: ss.GetColumns(),
			Values:  points,
		}
		res.Series = append(res.Series, row)
	}
	results := client090.Response{
		Results: []client090.Result{res},
		Err:     err,
	}

	return &results, err
}

func NewInfluxdbClient_v088(conf map[string]interface{}) (*InfluxdbClient_v088, error) {
	tmp_conf := map[string]interface{}{}
	for key, value := range conf {
		tmp_conf[strings.ToLower(key)] = value
	}

	host, ok := tmp_conf["host"]
	if ok != true || host == "" {
		return nil, fmt.Errorf("version 0.8.8, config missing: Host")
	}

	username, ok := tmp_conf["username"]
	if ok != true {
		return nil, fmt.Errorf("version 0.8.8, config missing: Username")
	}

	password, ok := tmp_conf["password"]
	if ok != true {
		return nil, fmt.Errorf("version 0.8.8, config missing: Password")
	}

	database, ok := tmp_conf["database"]
	if ok != true || database == "" {
		return nil, fmt.Errorf("version 0.8.8, config missing: Database")
	}

	flattemplate, ok := tmp_conf["flattemplate"]
	if ok != true {
		return nil, fmt.Errorf("version 0.8.8, config missing: FlatTemplate")
	}

	issecure, ok := tmp_conf["issecure"]
	if ok != true {
		issecure = false
	}

	c := &client088.ClientConfig{
		Host:     host.(string),
		Username: username.(string),
		Password: password.(string),
		Database: database.(string),
		IsSecure: issecure.(bool),
	}

	schema := "http"
	if c.IsSecure == true {
		schema = "https"
	}

	cli, err := client088.NewClient(c)
	if err != nil {
		return &InfluxdbClient_v088{}, fmt.Errorf("version 0.8.8, cannot create client: %v", err)
	}
	return &InfluxdbClient_v088{
		client:     cli,
		schema:     schema,
		username:   username.(string),
		password:   password.(string),
		host:       host.(string),
		httpClient: &http.Client{},
		flat_tpl:   flattemplate.(string),
	}, nil
}

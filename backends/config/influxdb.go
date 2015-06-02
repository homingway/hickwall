package config

type Transport_influxdb struct {
	// Enabled        bool

	Version        string
	Interval       string
	Max_batch_size int

	// Client Config
	Host     string // for v0.8.8
	URL      string // for v0.9.0
	Username string
	Password string
	Database string

	// Write Config
	RetentionPolicy string
	FlatTemplate    string

	Backfill_enabled              bool
	Backfill_interval             string
	Backfill_handsoff             bool
	Backfill_latency_threshold_ms int
	Backfill_cool_down            string

	Merge_Requests bool // try best to merge small group of points to no more than max_batch_size
}

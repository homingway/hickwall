package config

type Transport_file struct {
	Enabled        bool   `json:"enabled"`
	Flush_Interval string `json:"flush_interval"`
	Path           string `json:"path"`

	// TODO: max_size, max_rotation
	Max_size     int `json:"max_size"`
	Max_rotation int `json:"max_rotation"`
}

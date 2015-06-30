package config

type Transport_elasticsearch struct {

	// Client Config
	URL   string `json:"url"`
	Index string `json:"index"`
	Type  string `json:"type"`
	// Timeout int    `json:"timeout"`
}

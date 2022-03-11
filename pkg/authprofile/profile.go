package authprofile

type Profile struct {
	Name                string `json:"name,omitempty"`
	URL                 string `json:"url,omitempty"`
	Key                 string `json:"key,omitempty"`
	Secret              string `json:"secret,omitempty"`
	Username            string `json:"username,omitempty"`
	Password            string `json:"password,omitempty"`
	SkipServerCertValid bool   `json:"skip_server_cert_check,omitempty"`
}

type Profiles struct {

	// Default profile name, used as an index into the map

	Default  string             `json:"default,omitempty"`
	Profiles map[string]Profile `json:"profiles,omitempty"`
}

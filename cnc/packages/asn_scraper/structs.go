package asn_scraper

type ASN struct {
	Asn         string      `json:"asn"`
	Name        string      `json:"name"`
	Country     string      `json:"country"`
	Allocated   string      `json:"allocated"`
	Registry    string      `json:"registry"`
	Domain      string      `json:"domain"`
	NumIps      int         `json:"num_ips"`
	Type        string      `json:"type"`
	Prefixes    []Prefixes  `json:"prefixes"`
	Prefixes6   []Prefixes6 `json:"prefixes6"`
	Peers       []string    `json:"peers"`
	Upstreams   []string    `json:"upstreams"`
	Downstreams []string    `json:"downstreams"`
}

type Prefixes struct {
	Netblock string  `json:"netblock"`
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	Country  string  `json:"country"`
	Size     string  `json:"size"`
	Status   string  `json:"status"`
	Domain   *string `json:"domain"`
}

type Prefixes6 struct {
	Netblock string `json:"netblock"`
	Id       string `json:"id"`
	Name     string `json:"name"`
	Country  string `json:"country"`
	Size     string `json:"size"`
	Status   string `json:"status"`
	Domain   string `json:"domain"`
}

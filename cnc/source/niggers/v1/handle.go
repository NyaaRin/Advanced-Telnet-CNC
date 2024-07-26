package slave_v1

import (
	"net"
	"strings"
	"time"
)

var (
	FakeCounting = false
)

func Handle(conn net.Conn, source string) {
	client := &Client{
		Conn:   conn,
		Source: source,
	}

	if len(source) < 1 {
		client.Source = "undefined"
	}

	client.Source = strings.ToLower(client.Source)

	// Mapping of source strings to corresponding client sources
	mapping := map[string]string{
		/*"goahead":   "goahead",
		"dlink":     "undefined",
		"fiberrt":   "fiber",
		"avtech":    "avtech",
		"arm7":      "avtech",
		"c.":        "dvr",
		"dvr":       "dvr",
		"ipcam":     "dvr",
		"r190w":     "cnpilot",
		"iphone":    "china",
		"gponfiber": "boa",
		"boa":       "boa",
		"gpon":      "boa",
		"zhone":     "undefined",
		"weed":      "undefined",
		"baicells":  "undefined",
		"hongdian":  "undefined",
		"dm900":     "undefined",
		"lilin":     "dvr",
		"jaws":      "undefined",
		"gozy":      "undefined",
		"utt":       "undefined",
		"gargoyle":  "undefined",
		"ruijie":    "undefined",
		"brr":       "brazil",
		"spain":     "dvr",
		"cnr":       "boa",
		"asus":      "asus",
		"newzte":    "zte",
		"china":     "undefined",
		"wavlink":   "undefined",
		"webproc":   "undefined",
		"faith":     "undefined",
		"totolink":  "undefined",*/
	}

	for key, value := range mapping {
		if strings.Contains(source, key) {
			client.Source = value
		}
	}

	if len(client.Source) <= 1 {
		return
	}

	if client.Source == "e" {
		return
	}

	if client.Source == "teastick-telnet" {
		return
	}

	if client.Source == "teastick-lillin" {
		return
	}

	Add <- client
	defer func() {
		Delete <- client
	}()

	client.Handle()
}

func (c *Client) Handle() error {
	buf := make([]byte, 2)
	for {
		err := c.SetDeadline(time.Now().Add(180 * time.Second))
		if err != nil {
			return err
		}
		if n, err := c.Conn.Read(buf); err != nil || n != len(buf) {
			return err
		}
		if n, err := c.Conn.Write(buf); err != nil || n != len(buf) {
			return err
		}
	}
}

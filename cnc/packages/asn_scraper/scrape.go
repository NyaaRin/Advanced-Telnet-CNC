package asn_scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type Scraper struct {
	token string
}

type Cidr struct {
	Prefix  net.IP
	Netmask int
}

func NewScraper(token string) *Scraper {
	return &Scraper{token: token}
}

func (s *Scraper) CIDR(as string) (numOfIPs int, cidr []*Cidr, err error) {
	var cidrs = make([]*Cidr, 0)
	rwq, err := http.Get(fmt.Sprintf("http://ipinfo.io/%s?token=%s", as, s.token))
	if err != nil {
		return 0, nil, err
	}

	all, err := io.ReadAll(rwq.Body)
	if err != nil {
		return 0, nil, err
	}

	var asn *ASN
	err = json.Unmarshal(all, &asn)
	if err != nil {
		return 0, nil, err
	}

	for _, prefix := range asn.Prefixes {
		lol := strings.Split(prefix.Netblock, "/")
		atoi, err := strconv.Atoi(lol[1])
		if err != nil {
			return 0, nil, err
		}

		cidrs = append(cidrs, &Cidr{
			Prefix:  net.ParseIP(lol[0]),
			Netmask: atoi,
		})
	}

	if len(cidrs) > 255 {
		return asn.NumIps, cidrs[:255], err
	}

	return asn.NumIps, cidrs, err
}

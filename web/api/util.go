package api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func prettifyJSON(in string) (string, error) {
	var out bytes.Buffer

	err := json.Indent(&out, []byte(in), "", "  ")
	if err != nil {
		return "", err
	}

	return out.String(), nil
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	body, _ := prettifyJSON("{\"error\":\"" + message + "\"}")
	w.Write([]byte(body))
}

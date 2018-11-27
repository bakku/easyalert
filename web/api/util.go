package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
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
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)
	body, _ := prettifyJSON("{\"error\":\"" + message + "\"}")
	w.Write([]byte(body))
}

func getUserToken(r *http.Request) (string, bool) {
	rawToken := r.Header.Get("Authorization")
	if rawToken == "" {
		return "", false
	}

	splittedToken := strings.Split(rawToken, " ")
	if len(splittedToken) != 2 {
		return "", false
	}

	return splittedToken[1], true
}

package response

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// http response with json
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token, X-API-KEY")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

// decode the body data
func DecodeRequest(body io.ReadCloser, model interface{}) error {
	decoder := json.NewDecoder(body)
	return decoder.Decode(&model)
}

// http response with error [Bad Request]
func Error(w http.ResponseWriter) {
	RespondWithJSON(w, http.StatusBadRequest, "try again later")
}

func RespondWithHTML(w http.ResponseWriter, code int, html string) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
	// json.NewEncoder(w).Encode(payload)
	fmt.Fprintf(w, html)
}

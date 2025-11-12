package httpx

import (
	"encoding/json"
	"net/http"
)

func WriteJSON(writer http.ResponseWriter, code int, value any) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	_ = json.NewEncoder(writer).Encode(value)
}

func WriteErr(writer http.ResponseWriter, code int, err error) {
	WriteJSON(writer, code, map[string]string{"error": err.Error()})
}

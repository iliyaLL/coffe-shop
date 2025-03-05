package utils

import (
	"encoding/json"
	"net/http"
)

// sending responses in the json format
//
//	{
//		"error": "Internal Server Error"
//	}
type Response map[string]interface{}

func SendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

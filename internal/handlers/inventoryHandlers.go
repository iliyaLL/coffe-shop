package handlers

import (
	"net/http"
)

func (app *application) inventoryCreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello world"))
}

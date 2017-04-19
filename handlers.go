package main

import (
	"net/http"
)

func createOrderHandler(w http.ResponseWriter, r *http.Request)  {
    http.Error(w, "Not Implemented", http.StatusNotImplemented)
}
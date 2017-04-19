package orders

import (
    "net/http"

    "github.com/gorilla/mux"
)

func Routes() http.Handler {
    r := mux.NewRouter()
    r.HandleFunc("/orders", createOrderHandler).Methods("POST")

    return r
}

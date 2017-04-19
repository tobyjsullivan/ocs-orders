package orders

import (
    "net/http"
    "encoding/json"
    "github.com/satori/go.uuid"
)

type req struct {
    Name string `json:"name"`
    Phone string `json:"phone"`
    Address1 string `json:"address1"`
    Address2 string `json:"address2"`
    PostalCode string `json:"postalCode"`
    Instructions string `json:"additionalInstructions"`
}

type order struct {
    ID string `json:"id"`
    *req
}

func createOrderHandler(w http.ResponseWriter, r *http.Request)  {
    decoder := json.NewDecoder(r.Body)
    var request req
    err := decoder.Decode(&request)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
    }

    defer r.Body.Close()

    order := &order{
        ID: uuid.NewV4().String(),
        req: &request,
    }

    err = fireOrderAccepted(order)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }

    encoder := json.NewEncoder(w)
    encoder.Encode(order)
}

func fireOrderAccepted(order *order) error {
    // TODO Submit order to OrderAccepted queue

    return nil
}
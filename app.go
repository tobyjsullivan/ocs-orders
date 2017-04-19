package main

import (
    "os"

    "github.com/codegangsta/negroni"
    "github.com/tobyjsullivan/ocs-orders/orders"
)

func main() {
    n := negroni.New()
    n.UseHandler(orders.Routes())

    port := os.Getenv("PORT")
    if port == "" {
        port = "3000"
    }

    n.Run(":" + port)
}

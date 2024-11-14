package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
    "data-streaming/api"
)

func main() {
    router := mux.NewRouter()
    api.InitializeRoutes(router)

    log.Println("Starting server on :8080")
    log.Fatal(http.ListenAndServe(":8080", router))
}

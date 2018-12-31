package main

import (
    "net/http"

    "github.com/go-pg/pg"
    "github.com/gorilla/mux"
)

func main() {
    db := pg.Connect(&pg.Options{
        User:     "admingoforfun",
        Password: "admingoforfun",
    })

    mux := mux.NewRouter()
    server := newServer(db, mux)
    http.ListenAndServe(":8080", server)
}

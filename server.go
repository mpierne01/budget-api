// server.go
package main

type server struct {
    db *pg.DB
    mux *mux.Router
}
// server.go
package main


type server struct {
    db *pg.DB
    mux *mux.Router
}

func newServer(db *pg.DB, mux *mux.Router) *server {
    s := server{db, mux}
    s.routes() // register handlers
    return &s
}
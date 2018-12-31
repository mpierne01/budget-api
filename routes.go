// routes.go
package main

const BudgetPath = "/api/budgets"

func (s *server) routes() {
    s.mux.HandleFunc(BudgetPath, s.getBudgets).Methods("GET")
}

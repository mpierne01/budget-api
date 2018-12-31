// budgets_test.go
package main

import (
    "net/http"
    "os"
    "testing"

    "github.com/mpierne01/common/test"
    "github.com/go-pg/pg"
    "github.com/go-pg/pg/orm"
    "github.com/gorilla/mux"
    "github.com/pquerna/ffjson/ffjson"
    "github.com/stretchr/testify/assert"
)

var testServer *server
var badServer *server

func TestMain(m *testing.M) {
    testServer = newServer(
        pg.Connect(&pg.Options{
            User:     "admingoforfun",
            Password: "admingoforfun",
            Database: "admingoforfun_test",
        }),
        mux.NewRouter(),
    )

    badServer = newServer(
        pg.Connect(&pg.Options{
            User:     "not_found",
            Password: "goforfun",
            Database: "goforfun",
        }),
        mux.NewRouter(),
    )

        // Here we create a temporary table to store each test case
        // data and follow isolation which would be dropped after.
    testServer.db.CreateTable(&budget{}, &orm.CreateTableOptions{
        Temp: true,
    })

    os.Exit(m.Run())
}

func TestGetBudgets_EmptyResponse(t *testing.T) {
    var body []budget
    res, err := test.DoRequest(testServer, "GET", BudgetPath, nil)

    ffjson.NewDecoder().DecodeReader(res.Body, &body)
    assert.NoError(t, err)
    assert.Len(t, body, 0)
    assert.Equal(t, res.Code, http.StatusOK)
}

func TestGetBudgets_NormalResponse(t *testing.T) {
    var body []budget
    budgets, err := CreateBudgetListFactory(testServer.db, 10)
    assert.NoError(t, err)

    res, err := test.DoRequest(testServer, "GET", BudgetPath, nil)
    assert.Equal(t, http.StatusOK, res.Code)
    assert.NoError(t, err)

    ffjson.NewDecoder().DecodeReader(res.Body, &body)
    assert.Len(t, body, 10)
    assert.Equal(t, budgets, &body)
}

func TestGetBudgets_DatabaseError(t *testing.T) {
    var body Error
    res, err := test.DoRequest(badServer, "GET", BudgetPath, nil)

    ffjson.NewDecoder().DecodeReader(res.Body, &body)
    assert.NoError(t, err)
    assert.Equal(t, DatabaseError, &body)
    assert.Equal(t, http.StatusInternalServerError, res.Code)
}

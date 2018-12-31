// budgets_test.go
package main

import (
    "bytes"
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
            Database: "mpierne01_test",
        }),
        mux.NewRouter(),
    )

    badServer = newServer(
        pg.Connect(&pg.Options{
            User:     "not_found",
            Password: "admingoforfun",
            Database: "admingoforfun",
        }),
        mux.NewRouter(),
    )

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

func TestCreateBudget(t *testing.T) {
    var body budget
    byt, err := ffjson.Marshal(&budget{Amount: 1000.4})
    rdr := bytes.NewReader(byt)

    res, err := test.DoRequest(testServer, "POST", BudgetPath, rdr)

    ffjson.NewDecoder().DecodeReader(res.Body, &body)
    assert.NoError(t, err)
    assert.Equal(t, 1000.4, body.Amount)
    assert.Equal(t, http.StatusOK, res.Code)
}

func TestCreateBudget_BadParamError(t *testing.T) {
    var body Error
    res, err := test.DoRequest(testServer, "POST", BudgetPath,
        bytes.NewReader([]byte{}))

    ffjson.NewDecoder().DecodeReader(res.Body, &body)
    assert.NoError(t, err)
    assert.Equal(t, BadParamError, &body)
    assert.Equal(t, http.StatusBadRequest, res.Code)
}

func TestCreateBudget_DatabaseError(t *testing.T) {
    var body Error
    byt, err := ffjson.Marshal(&budget{Amount: 1000.4})
    rdr := bytes.NewReader(byt)

    res, err := test.DoRequest(badServer, "POST", BudgetPath, rdr)

    ffjson.NewDecoder().DecodeReader(res.Body, &body)
    assert.NoError(t, err)
    assert.Equal(t, DatabaseError, &body)
    assert.Equal(t, http.StatusInternalServerError, res.Code)
}

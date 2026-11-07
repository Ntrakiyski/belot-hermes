package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddBotsReturnsRequesterView(t *testing.T) {
	t.Setenv("ENABLE_PRACTICE_BOTS", "1")
	s := newStore()
	table := &table{ID: "TEST"}
	token := table.seatPlayer(0, "Tester")
	s.tables[table.ID] = table
	mux := http.NewServeMux()
	s.routes(mux)

	request := httptest.NewRequest("POST", "/api/tables/TEST/bots", nil)
	response := httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	if response.Code != http.StatusForbidden || table.Seats[1] != nil {
		t.Fatal("unauthorized request must not fill seats")
	}

	request = httptest.NewRequest("POST", "/api/tables/TEST/bots", nil)
	request.Header.Set("X-Token", token)
	response = httptest.NewRecorder()
	mux.ServeHTTP(response, request)
	var v view
	if err := json.Unmarshal(response.Body.Bytes(), &v); err != nil {
		t.Fatal(err)
	}
	if response.Code != http.StatusOK || v.Table != "TEST" || !v.Practice || !v.Started || len(v.Seats) != 4 || len(v.Hand) != 5 {
		t.Fatalf("expected full requester view after filling practice seats, got status=%d view=%+v", response.Code, v)
	}
}

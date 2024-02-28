package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/adalbertjnr/ws-person/types"
)

func HTTPServer(httpServerAddr string) error {
	http.HandleFunc("/aggregate", handleAggregate)
	return http.ListenAndServe(httpServerAddr, nil)
}

func handleAggregate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		returnResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		return
	}
	data := new(types.Person)
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		returnResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	fmt.Println("receiving data from http client", data)
}

func returnResponse(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(&v)
}

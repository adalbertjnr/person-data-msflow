package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/adalbertjnr/ws-person/types"
	"github.com/sirupsen/logrus"
)

func HTTPServer(httpServerAddr string) error {
	aggSvc := NewDataStore()
	http.HandleFunc("/aggregate", handleAggregate(aggSvc))
	http.HandleFunc("/invoice/{id}", handleDataGetterByID(aggSvc))
	fmt.Println("http aggregator server listening on port", httpServerAddr)
	return http.ListenAndServe(httpServerAddr, nil)
}

func handleAggregate(agg Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			returnResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		data := types.Person{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			returnResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		fmt.Println("receiving data from http client", data)
		if err := agg.Insert(data); err != nil {
			logrus.Error("error inserting data into aggregator database")
			return
		}
	}
}

func handleDataGetterByID(agg Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			returnResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
			return
		}
		personId, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			returnResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		foundPerson, err := agg.Get(personId)
		if err != nil {
			returnResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		returnResponse(w, http.StatusOK, foundPerson)
	}
}

func returnResponse(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(&v)
}

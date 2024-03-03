package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/adalbertjnr/ws-person/aggregator/client"
	"github.com/sirupsen/logrus"
)

const (
	httpGatewayAddr    = ":5000"
	httpAggregatorAddr = ":3001"
)

func main() {
	c := NewInvoice(client.NewHTTPClientEndpoint(httpAggregatorAddr))
	http.HandleFunc("/invoice", httpHandlerWrapper(c.handleGetById))
	fmt.Println("gateway http running on port", httpGatewayAddr)
	log.Fatal(http.ListenAndServe(httpGatewayAddr, nil))
}

type Invoice struct {
	client client.ClientPicker
}

func NewInvoice(c client.ClientPicker) *Invoice {
	return &Invoice{
		client: c,
	}
}

func (i *Invoice) handleGetById(w http.ResponseWriter, r *http.Request) error {
	id, err := strconv.Atoi(r.URL.Query().Get("person_id"))
	if err != nil {
		return fmt.Errorf("error converting the url from query to int %d", id)
	}
	p, err := i.client.GetPersonById(context.TODO(), id)
	if err != nil {
		logrus.Error(err)
	}
	return jsonResponse(w, http.StatusOK, p)
}

type apiFn func(w http.ResponseWriter, r *http.Request) error

func httpHandlerWrapper(fn apiFn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			jsonResponse(w, http.StatusInternalServerError, err.Error())
		}
		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"uri":        r.RequestURI,
				"remoteAddr": r.RemoteAddr,
				"time":       start,
				"took":       time.Since(start).Seconds(),
			})
		}(time.Now())

	}
}

func jsonResponse(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

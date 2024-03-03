package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/adalbertjnr/ws-person/types"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

type HTTPMetrics struct {
	aggregateCounter prometheus.Counter
	aggregateLatency prometheus.Histogram
}

func NewHTTPMetrics(n string) *HTTPMetrics {
	reqC := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: fmt.Sprintf("http_%s_request_counter", n),
		Name:      "aggregator",
	})
	reqL := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: fmt.Sprintf("http_%s_request_latency", n),
		Name:      "aggregator",
	})
	return &HTTPMetrics{
		aggregateCounter: reqC,
		aggregateLatency: reqL,
	}
}

func (h *HTTPMetrics) instrumentHandlerWrapper(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			h.aggregateLatency.Observe(time.Since(start).Seconds())
		}(time.Now())
		h.aggregateCounter.Inc()
		next(w, r)
	}
}

func HTTPServer(httpServerAddr string) error {
	var (
		aggSvc = NewDataStore()
		ag     = NewHTTPMetrics("agg_handler")
		in     = NewHTTPMetrics("inv_handler")
	)

	http.HandleFunc("/aggregate", ag.instrumentHandlerWrapper(handleAggregate(aggSvc)))
	http.HandleFunc("/invoice", in.instrumentHandlerWrapper(handleDataGetterByID(aggSvc)))
	http.Handle("/metrics", promhttp.Handler())

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
		personId, err := strconv.Atoi(r.URL.Query().Get("person_id"))
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

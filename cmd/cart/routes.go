// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for orders
// ----------------------------------------------------------------------------

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/benc-uk/dapr-store/pkg/auth"
	"github.com/benc-uk/dapr-store/pkg/models"
	"github.com/benc-uk/dapr-store/pkg/problem"

	"github.com/gorilla/mux"
)

//
// All routes we need should be registered here
//
func (api API) addRoutes(router *mux.Router) {
	router.HandleFunc("/submit", auth.AuthMiddleware(api.submitOrder)).Methods("POST")
}

//
// Create a new order
//
func (api API) submitOrder(resp http.ResponseWriter, req *http.Request) {
	cl, _ := strconv.Atoi(req.Header.Get("content-length"))
	if cl <= 0 {
		problem.Problem{"json-error", "Zero length body", 400, "Body is required", serviceName}.HttpSend(resp)
		return
	}

	orderID := makeID(5)
	order := models.Order{}
	err := json.NewDecoder(req.Body).Decode(&order)

	// Some basic validation and checking on what we've been posted
	if err != nil {
		problem.Problem{"json-error", "JSON decoding error", 400, err.Error(), serviceName}.HttpSend(resp)
		return
	}
	if order.Amount <= 0 || len(order.Items) == 0 {
		problem.Problem{"json-error", "Malformed orders JSON", 400, "Validation failed, check orders schema", serviceName}.HttpSend(resp)
		return
	}
	order.ID = orderID
	order.Status = models.OrderNew

	jsonPayload, err := json.Marshal(order)
	if err != nil {
		problem.Problem{"json-error", "Order JSON marshalling error", 500, err.Error(), serviceName}.HttpSend(resp)
		return
	}

	daprURL := fmt.Sprintf("http://localhost:%d/v1.0/publish/%s", daprPort, daprTopicName)
	daprResp, err := http.Post(daprURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil || (daprResp.StatusCode < 200 || daprResp.StatusCode > 299) {
		problem.SendDaprProblem(daprURL, resp, daprResp, err, serviceName)
		return
	}

	// Send the order back, but this time it will have an id
	resp.Header().Set("Content-Type", "application/json")
	resultBytes, _ := json.Marshal(order)
	resp.Write(resultBytes)
}

//
// Scummy but functional ID generator
//
func makeID(length int) string {
	id := ""
	possible := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < length; i++ {
		id += string(possible[rand.Intn(len(possible)-1)])
	}

	return id
}

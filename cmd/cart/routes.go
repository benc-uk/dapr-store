// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for cart
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
		problem.Send("Zero length body", "err://json-error", resp, problem.HTTP400, nil, serviceName)
		return
	}

	orderID := makeID(5)
	order := models.Order{}
	err := json.NewDecoder(req.Body).Decode(&order)

	// Some basic validation and checking on what we've been posted
	if err != nil {
		problem.Send("Malformed orders JSON", "err://json-decode", resp, problem.HTTP400, err, serviceName)
		return
	}
	if order.Amount <= 0 || len(order.Items) == 0 {
		problem.Send("Malformed orders JSON", "err://json-error", resp, problem.HTTP400, nil, serviceName)
		return
	}
	order.ID = orderID
	order.Status = models.OrderNew

	jsonPayload, err := json.Marshal(order)
	if err != nil {
		problem.Send("Malformed orders JSON", "err://json-marshal", resp, nil, err, serviceName)
		return
	}

	daprURL := fmt.Sprintf("http://localhost:%d/v1.0/publish/%s", daprPort, daprTopicName)
	daprResp, err := http.Post(daprURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		problem.Send("Error publishing event", daprURL, resp, daprResp, err, serviceName)
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

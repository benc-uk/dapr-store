// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for cart
// ----------------------------------------------------------------------------

package main

import (
	"encoding/json"
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
		problem.New("err://body-missing", "Zero length body", 400, "Zero length body", daprHelper.AppInstanceName).Send(resp)
		return
	}

	orderID := makeID(5)
	order := models.Order{}
	err := json.NewDecoder(req.Body).Decode(&order)

	// Some basic validation and checking on what we've been posted
	if err != nil {
		problem.New("err://json-decode", "Malformed order JSON", 400, "JSON could not be decoded", daprHelper.AppInstanceName).Send(resp)
		return
	}
	if order.Amount <= 0 || len(order.Items) == 0 || order.Title == "" {
		problem.New("err://json-error", "Malformed order JSON", 400, "Order failed validation, check spec", daprHelper.AppInstanceName).Send(resp)
		return
	}
	order.ID = orderID
	order.Status = models.OrderNew

	prob := daprHelper.PublishMessage(ordersTopicName, order)
	if prob != nil {
		prob.Send(resp)
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

// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for orders
// ----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/benc-uk/dapr-store/cmd/orders/dapr"
	"github.com/benc-uk/dapr-store/cmd/orders/spec"
	"github.com/benc-uk/dapr-store/pkg/auth"
	"github.com/benc-uk/dapr-store/pkg/problem"
	"github.com/gorilla/mux"
)

//
// All routes we need should be registered here
//
func (api API) addRoutes(router *mux.Router) {

	router.HandleFunc("/dapr/subscribe", api.subscribeTopic)
	router.HandleFunc("/"+api.getTopicName(), api.receiveOrders)
	router.HandleFunc("/get/{id}", auth.AuthMiddleware(api.getOrder))
	router.HandleFunc("/getForUser/{username}", auth.AuthMiddleware(api.getOrdersForUser))
}

//
// Dapr pub-sub subscription endpoint, https://github.com/dapr/docs/blob/master/reference/api/pubsub_api.md
//
func (api API) subscribeTopic(resp http.ResponseWriter, req *http.Request) {
	// A simple JSON array of strings, each being a topic we subscribe to
	topicListJSON := fmt.Sprintf(`["%s"]`, api.getTopicName())
	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte(topicListJSON))
}

func (api API) getTopicName() string {
	topicName := api.service.(*dapr.OrderService).TopicName
	if topicName == "" {
		topicName = "orders-queue"
	}
	return topicName
}

//
// Dapr pub-sub topic receiver, https://github.com/dapr/docs/blob/master/reference/api/pubsub_api.md
//
func (api API) receiveOrders(resp http.ResponseWriter, req *http.Request) {
	type cloudevent struct {
		ID   string     `json:"id"`
		Data spec.Order `json:"data"`
	}
	event := &cloudevent{}

	err := json.NewDecoder(req.Body).Decode(&event)
	if err != nil {
		// Returning a non-200 will reschedule the received message
		problem.New("err://json-decode", "Event JSON decoding error", 500, err.Error(), api.ServiceName).Send(resp)
		return
	}
	log.Println("### Received event from orders pub/sub topic")

	// Save order in state with received status
	order := event.Data
	order.Status = spec.OrderReceived
	api.service.AddOrder(order)
	prob := api.service.AddOrder(order) //api.SaveState(daprStoreName, order.ID, order)
	if prob != nil {
		// Returning a non-200 will reschedule the received message
		prob := err.(*problem.Problem)
		prob.Send(resp)
		return
	}
	log.Printf("### Order %s was saved to state store\n", order.ID)

	// Save order to blob storage as a text file "report"
	// Also email to the user via SendGrid
	// For these to work configure the components in cmd/orders/components
	// If un-configured then nothing happens, and no output is send or generated
	api.service.EmailNotify(order)
	api.service.SaveReport(order)

	// Fake background order processing
	time.AfterFunc(30*time.Second, func() {
		api.service.SetStatus(&order, spec.OrderProcessing)
	})

	// Fake background order completion
	time.AfterFunc(120*time.Second, func() {
		api.service.SetStatus(&order, spec.OrderComplete)
	})

	// Send success
	resp.WriteHeader(200)
}

//
// Fetch existing order by id
//
func (api API) getOrder(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	order, err := api.service.GetOrder(vars["id"])
	if err != nil {
		prob := err.(*problem.Problem)
		prob.Send(resp)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(order)
	resp.Write(json)
}

//
// Fetch all orders for a given user
//
func (api API) getOrdersForUser(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	orders, err := api.service.GetOrdersForUser(vars["username"])
	if err != nil {
		prob := err.(*problem.Problem)
		prob.Send(resp)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	json, _ := json.Marshal(orders)
	resp.Write(json)
}

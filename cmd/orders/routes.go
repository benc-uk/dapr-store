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

	"github.com/benc-uk/dapr-store/pkg/auth"
	"github.com/benc-uk/dapr-store/pkg/models"
	"github.com/benc-uk/dapr-store/pkg/problem"
	"github.com/gorilla/mux"
)

//
// All routes we need should be registered here
//
func (api API) addRoutes(router *mux.Router) {
	router.HandleFunc("/dapr/subscribe", api.subscribeTopic)
	router.HandleFunc("/"+daprHelper.PubsubQueueName, api.receiveOrders)
	router.HandleFunc("/get/{id}", auth.AuthMiddleware(api.getOrder))
	router.HandleFunc("/getForUser/{username}", auth.AuthMiddleware(api.getOrdersForUser))
}

//
// Dapr pub-sub subscription endpoint, https://github.com/dapr/docs/blob/master/reference/api/pubsub_api.md
//
func (api API) subscribeTopic(resp http.ResponseWriter, req *http.Request) {
	// A simple JSON array of strings, each being a topic we subscribe to
	topicListJSON := fmt.Sprintf(`["%s"]`, daprHelper.PubsubQueueName)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte(topicListJSON))
}

//
// Dapr pub-sub topic receiver, https://github.com/dapr/docs/blob/master/reference/api/pubsub_api.md
//
func (api API) receiveOrders(resp http.ResponseWriter, req *http.Request) {
	type cloudevent struct {
		ID   string       `json:"id"`
		Data models.Order `json:"data"`
	}
	event := &cloudevent{}

	err := json.NewDecoder(req.Body).Decode(&event)
	if err != nil {
		// Returning a non-200 will reschedule the received message
		problem.New("err://json-decode", "Event JSON decoding error", 500, err.Error(), daprHelper.AppInstanceName).Send(resp)
		return
	}
	log.Printf("### Received event from pub/sub topic: %s\n", daprHelper.PubsubQueueName)

	// Save order in state with received status
	order := event.Data
	order.Status = models.OrderReceived
	prob := daprHelper.SaveState(order.ID, order)
	if prob != nil {
		// Returning a non-200 will reschedule the received message
		prob.Send(resp)
		return
	}
	log.Printf("### Order %s was saved to state store\n", order.ID)

	// Save order to blob storage as a text file "notification"
	outputMetadata := map[string]string{
		"ContentType": "text/plain",
		"blobName":    "order_" + order.ID + ".txt",
	}
	outputData := "----------\nOrder title:" + order.Title + "\nOrder ID: " + order.ID + "\nUser: " + order.ForUser + "\nAmount: " + fmt.Sprintf("%f", order.Amount) + "\n----------"
	prob = daprHelper.SendOutput("orders-notify", outputData, outputMetadata)
	if prob != nil {
		log.Printf("### Problem sending to notification output: %+v", prob)
	}

	// Now create or update the user's orders index, which is keyed on their username
	// And is simply an array of orderIds (strings)
	userOrders := []string{}
	// !NOTE! We use the username as a key in the orders state set, to hold an index of orders
	data, prob := daprHelper.GetState(order.ForUser)
	// Ignore any problem, it's possible it doesn't exist yet (user's first order)
	err = json.Unmarshal(data, &userOrders)

	alreadyExists := false
	log.Printf("### userOrders is %v", userOrders)
	for _, oid := range userOrders {
		if order.ID == oid {
			alreadyExists = true
		}
	}

	if !alreadyExists {
		userOrders = append(userOrders, order.ID)
	} else {
		log.Printf("### Warning, duplicate order '%s' for user '%s' detected", order.ID, order.ForUser)
	}

	// Save new list back
	prob = daprHelper.SaveState(order.ForUser, userOrders)
	if prob != nil {
		// Returning a non-200 will reschedule the received message
		prob.Send(resp)
		log.Printf("### Error!, unable to save order list for user '%s'", order.ForUser)
	}

	// Fake background order processing
	time.AfterFunc(30*time.Second, func() {
		order.Status = models.OrderProcessing
		prob = daprHelper.SaveState(order.ID, order)
		if prob != nil {
			log.Printf("### Warning, order processing failed: %s", prob.Error())
		} else {
			log.Printf("### Order %s was moved to status: %s", order.ID, order.Status)
		}
	})

	// Fake background order completion
	time.AfterFunc(120*time.Second, func() {
		order.Status = models.OrderComplete
		prob = daprHelper.SaveState(order.ID, order)
		if prob != nil {
			log.Printf("### Warning, order completion failed: %s", prob.Error())
		} else {
			log.Printf("### Order %s was moved to status: %s", order.ID, order.Status)
		}
	})

	// Send success
	resp.WriteHeader(200)
}

//
// Fetch existing order by id
//
func (api API) getOrder(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	data, prob := daprHelper.GetState(vars["id"])
	if prob != nil {
		prob.Send(resp)
		return
	}
	if len(data) <= 0 {
		problem.New("err://not-found", "No data returned", 404, "Order id:'"+vars["id"]+"' not found", daprHelper.AppInstanceName).Send(resp)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(data)
}

//
// Fetch all orders for a given user
//
func (api API) getOrdersForUser(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	// !NOTE! We use the username as a key in the orders state set, to hold an index of orders
	data, prob := daprHelper.GetState(vars["username"])
	if prob != nil {
		prob.Send(resp)
		return
	}

	// If no data, just return an empty JSON array
	if len(data) <= 0 {
		resp.Header().Set("Content-Type", "application/json")
		resp.Write([]byte(`[]`))
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(data)
}

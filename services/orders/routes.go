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
	"net/http"
	"time"

	"github.com/benc-uk/dapr-store/common"
	"github.com/gorilla/mux"
)

//
// All routes we need should be registered here
//
func (api API) addRoutes(router *mux.Router) {
	router.HandleFunc("/dapr/subscribe", api.subscribeTopic)
	router.HandleFunc("/"+daprTopicName, api.receiveOrders)
	router.HandleFunc("/get/{id}", common.AuthMiddleware(api.getOrder))
}

//
// Dapr pub-sub subscription endpoint, https://github.com/dapr/docs/blob/master/reference/api/pubsub_api.md
//
func (api API) subscribeTopic(resp http.ResponseWriter, req *http.Request) {
	// A simple JSON array of strings, each being a topic we subscribe to
	topicListJSON := fmt.Sprintf(`["%s"]`, daprTopicName)
	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte(topicListJSON))
}

//
// Dapr pub-sub topic receiver, https://github.com/dapr/docs/blob/master/reference/api/pubsub_api.md
//
func (api API) receiveOrders(resp http.ResponseWriter, req *http.Request) {
	type cloudevent struct {
		ID   string       `json:"id"`
		Data common.Order `json:"data"`
	}
	event := &cloudevent{}

	err := json.NewDecoder(req.Body).Decode(&event)
	if err != nil {
		common.Problem{"json-error", "JSON decoding error", 500, err.Error(), serviceName}.HttpSend(resp)
		return
	}

	// Save order in state with received status
	order := event.Data
	order.Status = common.OrderReceived
	err = common.SaveState(resp, daprPort, daprStoreName, serviceName, order.ID, order)
	if err != nil {
		return // Error will have already been written to resp
	}

	// Update the user adding the orderId to their list of owned orders
	client := &http.Client{}
	addorderPutReq, err := http.NewRequest("PUT", fmt.Sprintf(`http://localhost:%d/v1.0/invoke/users/method/addorder/%s/%s`, daprPort, order.ForUser, order.ID), nil)
	if err == nil {
		_, err := client.Do(addorderPutReq)
		if err != nil {
			resp.WriteHeader(500)
			return
		}
	}

	// Fake background order processing
	time.AfterFunc(30*time.Second, func() {
		order.Status = common.OrderProcessing
		err = common.SaveState(resp, daprPort, daprStoreName, serviceName, order.ID, order)
	})

	// Fake background order completion
	time.AfterFunc(120*time.Second, func() {
		order.Status = common.OrderComplete
		err = common.SaveState(resp, daprPort, daprStoreName, serviceName, order.ID, order)
	})

	// Send success
	resp.WriteHeader(200)
}

//
// Fetch existing order by id
//
func (api API) getOrder(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	data, err := common.GetState(resp, daprPort, daprStoreName, serviceName, vars["id"])

	if len(data) <= 0 {
		common.Problem{"order-not-found", vars["id"] + " not found", 404, "Order does not exist", serviceName}.HttpSend(resp)
		return
	}

	if err != nil {
		return // Error will have already been written to resp
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(data)
}

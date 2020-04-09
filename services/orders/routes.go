package main

//
// Basic REST API microservice, template/reference code
// Ben Coleman, July 2019, v1
//

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/benc-uk/dapr-store/common"
	"github.com/gorilla/mux"
)

// All routes we need should be here
func (api API) addRoutes(router *mux.Router) {
	router.HandleFunc("/new", api.newOrder).Methods("POST")
	router.HandleFunc("/get/{id}", api.getOrder)
}

//
// Create a new order
//
func (api API) newOrder(resp http.ResponseWriter, req *http.Request) {
	cl, _ := strconv.Atoi(req.Header.Get("content-length"))
	if cl <= 0 {
		common.Problem{"json-error", "Zero length body", 400, "Body is required", serviceName}.HttpSend(resp)
		return
	}

	orderID := makeID(5)
	order := Order{}
	err := json.NewDecoder(req.Body).Decode(&order)

	// Some basic validation and checking on what we've been posted
	if err != nil {
		common.Problem{"json-error", "JSON decoding error", 400, err.Error(), serviceName}.HttpSend(resp)
		return
	}
	if order.Amount <= 0 || len(order.ProductID) == 0 {
		common.Problem{"json-error", "Malformed orders JSON", 400, "Validation failed, check orders schema", serviceName}.HttpSend(resp)
		return
	}
	order.ID = orderID

	err = common.SaveState(resp, daprPort, daprStateStore, serviceName, orderID, order)
	if err != nil {
		return // Error will have already been written to resp
	}

	// Send the order back, but this time it will have an id
	resp.Header().Set("Content-Type", "application/json")
	resultBytes, _ := json.Marshal(order)
	resp.Write(resultBytes)
}

//
// Fetch existing order
//
func (api API) getOrder(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	data, err := common.GetState(resp, daprPort, daprStateStore, serviceName, vars["id"])

	if err != nil {
		return // Error will have already been written to resp
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(data)
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

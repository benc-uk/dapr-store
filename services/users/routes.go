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
	"strconv"

	"github.com/benc-uk/dapr-store/common"
	"github.com/gorilla/mux"
)

//
// All routes we need should be registered here
//
func (api API) addRoutes(router *mux.Router) {
	router.HandleFunc("/register", api.registerUser).Methods("POST")
	router.HandleFunc("/get/{username}", api.getUser)
	router.HandleFunc("/addorder/{username}/{orderId}", api.addOrderToUser).Methods("PUT")
}

//
// Register new user
//
func (api API) registerUser(resp http.ResponseWriter, req *http.Request) {

	cl, _ := strconv.Atoi(req.Header.Get("content-length"))
	if cl <= 0 {
		common.Problem{"json-error", "Zero length body", 400, "Body is required", serviceName}.HttpSend(resp)
		return
	}

	user := common.User{}
	err := json.NewDecoder(req.Body).Decode(&user)

	// Some basic validation and checking on what we've been posted
	if err != nil {
		common.Problem{"json-error", "JSON decoding error", 400, err.Error(), serviceName}.HttpSend(resp)
		return
	}
	if len(user.DisplayName) == 0 || len(user.Username) == 0 {
		common.Problem{"json-error", "Malformed user JSON", 400, "Validation failed, check user schema", serviceName}.HttpSend(resp)
		return
	}
	log.Printf("### Registering user %+v\n", user)

	// Check is user already registered
	data, err := common.GetState(resp, daprPort, daprStoreName, serviceName, user.Username)
	if err != nil {
		return // Error will have already been written to resp
	}
	log.Printf("### Existing user data %+v\n", string(data))
	if len(data) > 0 {
		common.Problem{"user-exists", user.Username + " already registered", 400, "User is already registered!", serviceName}.HttpSend(resp)
		return
	}

	err = common.SaveState(resp, daprPort, daprStoreName, serviceName, user.Username, user)
	if err != nil {
		return // Error will have already been written to resp
	}

	// Send success message back
	resp.Header().Set("Content-Type", "application/json")
	resp.Write([]byte(fmt.Sprintf(`{"registrationStatus":"success", "username":"%s"`, user.Username)))
}

//
// Fetch existing user
//
func (api API) getUser(resp http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	data, err := common.GetState(resp, daprPort, daprStoreName, serviceName, vars["username"])
	if err != nil {
		return // Error will have already been written to resp
	}

	if len(data) <= 0 {
		common.Problem{"user-not-found", vars["username"] + " not found", 404, "User is not registered", serviceName}.HttpSend(resp)
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.Write(data)
}

//
// Add orderId to user - !TODO! secure this so only callable from inside
//
func (api API) addOrderToUser(resp http.ResponseWriter, req *http.Request) {
	fmt.Printf("############################## HHHHHHHHHH\n\n\n")
	vars := mux.Vars(req)
	data, err := common.GetState(resp, daprPort, daprStoreName, serviceName, vars["username"])
	if err != nil {
		return // Error will have already been written to resp
	}

	if len(data) <= 0 {
		common.Problem{"user-not-found", vars["username"] + " not found", 404, "User is not registered", serviceName}.HttpSend(resp)
		return
	}

	user := common.User{}
	err = json.Unmarshal(data, &user)
	fmt.Printf("#### orderId %+v\n", vars["orderId"])
	fmt.Printf("#### username %+v\n", vars["username"])
	fmt.Printf("#### addOrderToUser %+v\n", user)
	orderID := vars["orderId"]
	alreadyExists := false
	for _, oid := range user.Orders {
		if orderID == oid {
			alreadyExists = true
		}
	}

	if !alreadyExists {
		user.Orders = append(user.Orders, orderID)
	} else {
		common.Problem{"order-exists", "No duplicate orders", 400, "Order '" + orderID + "' already assigned to user", serviceName}.HttpSend(resp)
		return
	}

	err = common.SaveState(resp, daprPort, daprStoreName, serviceName, vars["username"], user)
	if err != nil {
		return // Error will have already been written to resp
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(200)
}

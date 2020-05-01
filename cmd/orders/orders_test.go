// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Main set of tests for orders service and API + OrderService biz logic
// ----------------------------------------------------------------------------

package main

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/benc-uk/dapr-store/cmd/orders/mock"
	"github.com/benc-uk/dapr-store/cmd/orders/spec"
	"github.com/benc-uk/dapr-store/pkg/api"
	"github.com/benc-uk/dapr-store/pkg/apitests"
	"github.com/gorilla/mux"
)

func TestOrders(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	// Mock of CartService
	mockOrdersSvc := &mock.OrderService{}

	router := mux.NewRouter()
	api := API{
		api.NewBase("orders", "ignore", "ignore", true, router),
		mockOrdersSvc,
	}
	api.addRoutes(router)

	apitests.Run(t, router, testCases)

	// Rest of tests don't go through the router/api

	emptyOrder := spec.Order{}
	err := mockOrdersSvc.ProcessOrder(emptyOrder)
	if err != nil && strings.Contains(err.Error(), "validation") {
		t.Log("'process empty order' passed")
	} else {
		t.Error("'process empty order' failed")
	}

	badOrder := spec.Order{
		Items:   []string{"3"},
		ForUser: "",
	}
	err = mockOrdersSvc.ProcessOrder(badOrder)
	if err != nil && strings.Contains(err.Error(), "validation") {
		t.Log("'process invalid order' passed")
	} else {
		t.Error("'process invalid order' failed")
	}

	goodOrder := spec.Order{
		Items:   []string{"3"},
		ForUser: "test@example.net",
		Amount:  23.66,
		Status:  spec.OrderNew,
		Title:   "Test order",
		ID:      "test-01",
	}
	err = mockOrdersSvc.ProcessOrder(goodOrder)
	if err != nil {
		t.Errorf("'process valid new order' failed %+v", err)
	} else {
		t.Log("'process valid new order' passed")
	}

	newOrder, err := mockOrdersSvc.GetOrder("test-01")
	if err != nil {
		t.Errorf("'get new order' failed: %+v", err)
	} else {
		if newOrder.Status != spec.OrderReceived {
			t.Error("'get new order' failed")
		} else {
			t.Log("'get new order' passed")
		}
	}

	time.Sleep(time.Second * 3)
	newOrder, err = mockOrdersSvc.GetOrder("test-01")
	if err != nil {
		t.Errorf("'order processing completed' failed: %+v", err)
	} else {
		if newOrder.Status != spec.OrderComplete {
			t.Error("'order processing completed' failed")
		} else {
			t.Log("'order processing completed' passed")
		}
	}
}

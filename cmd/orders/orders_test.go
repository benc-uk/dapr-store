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

	t.Run("process empty order", func(t *testing.T) {
		emptyOrder := spec.Order{}
		err := mockOrdersSvc.ProcessOrder(emptyOrder)
		if err != nil && strings.Contains(err.Error(), "validation") {
		} else {
			t.Error("'process empty order' failed")
		}
	})

	t.Run("process invalid order", func(t *testing.T) {
		badOrder := spec.Order{
			LineItems: []spec.LineItem{},
			ForUser:   "",
		}
		err := mockOrdersSvc.ProcessOrder(badOrder)
		if err != nil && strings.Contains(err.Error(), "validation") {
		} else {
			t.Error("'process invalid order' failed")
		}
	})

	t.Run("process valid new order", func(t *testing.T) {
		goodOrder := mock.MockOrders[0]
		err := mockOrdersSvc.ProcessOrder(goodOrder)
		if err != nil {
			t.Errorf("'process valid new order' failed %+v", err)
		}
	})

	t.Run("get new order", func(t *testing.T) {
		newOrder, err := mockOrdersSvc.GetOrder("ord-mock")
		if err != nil {
			t.Errorf("'get new order' failed: %+v", err)
		} else {
			if newOrder.Status != spec.OrderReceived {
				t.Error("'get new order' failed")
			}
		}
	})

	t.Run("order processing completed", func(t *testing.T) {
		time.Sleep(time.Second * 3)
		newOrder, err := mockOrdersSvc.GetOrder("ord-mock")
		if err != nil {
			t.Errorf("'order processing completed' failed: %+v", err)
		} else {
			if newOrder.Status != spec.OrderComplete {
				t.Error("'order processing completed' failed")
			}
		}
	})
}

var testCases = []apitests.Test{
	{
		"get an existing order",
		"/get/ord-mock",
		"GET",
		"",
		"ord-mock",
		1,
		200,
	},
	{
		"get an non-existent order",
		"/get/foo",
		"GET",
		``,
		"not found", 1,
		404,
	},
}

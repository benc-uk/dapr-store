// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for orders
// ----------------------------------------------------------------------------

package main

import (
	"net/http"

	"github.com/benc-uk/dapr-store/cmd/orders/impl"
	"github.com/benc-uk/go-rest-api/pkg/auth"
	"github.com/benc-uk/go-rest-api/pkg/problem"
	"github.com/go-chi/chi/v5"
)

// All routes we need should be registered here
func (api API) addRoutes(router chi.Router, v auth.Validator) {
	router.Get("/get/{id}", v.Protect(api.getOrder))
	router.Get("/getForUser/{username}", v.Protect(api.getOrdersForUser))
}

// Fetch existing order by id
func (api API) getOrder(resp http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")

	order, err := api.service.GetOrder(id)
	if err != nil {
		if orderError, ok := err.(impl.OrdersError); ok && orderError.Error() == impl.NotFoundError {
			problem.Wrap(404, req.RequestURI, id, err).Send(resp)

			return
		}

		problem.Wrap(500, req.RequestURI, id, err).Send(resp)

		return
	}

	api.ReturnJSON(resp, order)
}

// Fetch all orders for a given user
func (api API) getOrdersForUser(resp http.ResponseWriter, req *http.Request) {
	username := chi.URLParam(req, "username")

	orders, err := api.service.GetOrdersForUser(username)
	if err != nil {
		problem.Wrap(500, req.RequestURI, username, err).Send(resp)

		return
	}

	api.ReturnJSON(resp, orders)
}

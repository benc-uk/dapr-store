// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Dapr compatible REST API service for cart
// ----------------------------------------------------------------------------

package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/benc-uk/dapr-store/cmd/cart/impl"
	"github.com/benc-uk/go-rest-api/pkg/auth"
	"github.com/benc-uk/go-rest-api/pkg/problem"
	"github.com/go-chi/chi/v5"
)

// All routes we need should be registered here
func (api API) addRoutes(router chi.Router, v auth.Validator) {
	router.Put("/setProduct/{userId}/{productId}/{count}", v.Protect(api.setProductCount))
	router.Get("/get/{userId}", v.Protect(api.getCart))
	router.Post("/submit", v.Protect(api.submitCart))
	router.Put("/clear/{userId}", v.Protect(api.clearCart))
}

func (api API) setProductCount(resp http.ResponseWriter, req *http.Request) {
	userID := chi.URLParam(req, "userId")
	productID := chi.URLParam(req, "productId")
	countString := chi.URLParam(req, "count")

	cart, err := api.service.Get(userID)
	if err != nil {
		problem.Wrap(500, req.RequestURI, userID, err).Send(resp)

		return
	}

	count, err := strconv.Atoi(countString)
	if err != nil {
		problem.Wrap(400, req.RequestURI, productID, err).Send(resp)
		return
	}

	err = api.service.SetProductCount(cart, productID, count)
	if err != nil {
		problem.Wrap(500, req.RequestURI, productID, err).Send(resp)

		return
	}

	resp.Header().Set("Content-Type", "application/json")

	api.ReturnJSON(resp, cart)
}

func (api API) getCart(resp http.ResponseWriter, req *http.Request) {
	userID := chi.URLParam(req, "userId")

	cart, err := api.service.Get(userID)

	if err != nil {
		problem.Wrap(500, req.RequestURI, userID, err).Send(resp)

		return
	}

	api.ReturnJSON(resp, cart)
}

func (api API) clearCart(resp http.ResponseWriter, req *http.Request) {
	userID := chi.URLParam(req, "userId")

	cart, err := api.service.Get(userID)
	if err != nil {
		problem.Wrap(500, req.RequestURI, userID, err).Send(resp)

		return
	}

	err = api.service.Clear(cart)
	if err != nil {
		log.Printf("### Warning failed to clear cart %s", err)
	}

	api.ReturnJSON(resp, cart)
}

func (api API) submitCart(resp http.ResponseWriter, req *http.Request) {
	userID := ""

	err := json.NewDecoder(req.Body).Decode(&userID)
	if err != nil {
		problem.Wrap(400, req.RequestURI, "none", err).Send(resp)
		return
	}

	if userID == "" {
		problem.Wrap(400, req.RequestURI, "none", errors.New("userId missing from request")).Send(resp)
		return
	}

	cart, err := api.service.Get(userID)
	if err != nil {
		problem.Wrap(500, req.RequestURI, userID, err).Send(resp)

		return
	}

	order, err := api.service.Submit(*cart)
	if err != nil {
		if cartErr, ok := err.(impl.CartError); ok && cartErr.Error() == impl.EmptyError {
			problem.Wrap(400, req.RequestURI, userID, cartErr).Send(resp)
			return
		}

		problem.Wrap(500, req.RequestURI, userID, err).Send(resp)

		return
	}

	// Send the _order_ back, created from submitting the cart
	api.ReturnJSON(resp, order)
}

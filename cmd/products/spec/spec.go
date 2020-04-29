// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Specification of the Product entity and service
// ----------------------------------------------------------------------------

package spec

// Product holds product data
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Cost        float32 `json:"cost"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	OnOffer     bool    `json:"onOffer"`
}

// UserService defines core CRUD methods a user service should have
type ProductService interface {
	SearchProducts(query string) ([]Product, error)
	QueryProducts(field, term string) ([]Product, error)
	AllProducts() ([]Product, error)
}

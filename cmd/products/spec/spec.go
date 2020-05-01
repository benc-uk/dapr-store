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

// ProductService defines core CRUD methods a products service should have
type ProductService interface {
	SearchProducts(string) ([]Product, error)
	QueryProducts(string, string) ([]Product, error)
	AllProducts() ([]Product, error)
}

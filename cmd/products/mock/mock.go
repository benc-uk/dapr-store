// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// Mock implementation of the ProductsService
// ----------------------------------------------------------------------------

package mock

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/benc-uk/dapr-store/cmd/products/spec"
)

// ProductsService mock version
type ProductsService struct {
}

// Load mock data
var mockProducts []spec.Product

func init() {
	mockJSON, err := ioutil.ReadFile("../../etc/mock-data/products.json")
	if err != nil {
		panic(err)
	}
	json.Unmarshal(mockJSON, &mockProducts)
}

// var products = []spec.Product{
// 	{
// 		ID:          "1",
// 		Name:        "A big hat",
// 		Cost:        23.99,
// 		Description: "It's really big, and green too",
// 		Image:       "massive-hat.jpg",
// 		OnOffer:     true,
// 	},
// 	{
// 		ID:          "2",
// 		Name:        "A small hat",
// 		Cost:        8.75,
// 		Description: "It's tiny, and pink too",
// 		Image:       "miniscule-hat.jpg",
// 		OnOffer:     false,
// 	},
// 	{
// 		ID:          "3",
// 		Name:        "A medium hat",
// 		Cost:        12.33,
// 		Description: "It's average size, and bright orange",
// 		Image:       "normal-hat.jpg",
// 		OnOffer:     false,
// 	},
// }

// SearchProducts mock/fake DB
func (s ProductsService) SearchProducts(query string) ([]spec.Product, error) {
	results := []spec.Product{}
	for _, prod := range mockProducts {
		if strings.Contains(prod.Name, query) || strings.Contains(prod.Description, query) {
			results = append(results, prod)
		}
	}
	return results, nil
}

// QueryProducts mock/fake DB
func (s ProductsService) QueryProducts(field, term string) ([]spec.Product, error) {
	results := []spec.Product{}
	for _, prod := range mockProducts {
		if field == "ID" && prod.ID == term {
			results = append(results, prod)
		}
		if field == "onoffer" {
			termBool, _ := strconv.ParseBool(term)
			if prod.OnOffer == termBool {
				results = append(results, prod)
			}
		}
	}
	return results, nil
}

// AllProducts mock/fake DB
func (s ProductsService) AllProducts() ([]spec.Product, error) {
	return mockProducts, nil
}

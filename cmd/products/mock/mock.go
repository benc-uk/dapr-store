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
type ProductService struct {
}

// Load mock data
var mockProducts []spec.Product

func init() {
	mockJSON, err := ioutil.ReadFile("../../testing/mock-data/products.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(mockJSON, &mockProducts)
	if err != nil {
		panic(err)
	}
}

// SearchProducts mock/fake DB
func (s ProductService) SearchProducts(query string) ([]spec.Product, error) {
	results := []spec.Product{}
	for _, prod := range mockProducts {
		if strings.Contains(prod.Name, query) || strings.Contains(prod.Description, query) {
			results = append(results, prod)
		}
	}
	return results, nil
}

// QueryProducts mock/fake DB
func (s ProductService) QueryProducts(field, term string) ([]spec.Product, error) {
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
func (s ProductService) AllProducts() ([]spec.Product, error) {
	return mockProducts, nil
}

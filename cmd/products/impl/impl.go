// ----------------------------------------------------------------------------
// Copyright (c) Ben Coleman, 2020
// Licensed under the MIT License.
//
// SQLite implementation of the ProductService
// ----------------------------------------------------------------------------

package impl

import (
	"database/sql"
	"log"

	"github.com/benc-uk/dapr-store/cmd/products/spec"
	"github.com/benc-uk/dapr-store/pkg/problem"
)

// ProductService is a Dapr based implementation of ProductService interface
type ProductService struct {
	*sql.DB
	serviceName string
}

// NewService creates a new ProductService
func NewService(serviceName string) *ProductService {
	db, err := sql.Open("sqlite3", "./sqlite.db")
	if err != nil {
		log.Panicf("### Failed to open database! %+v\n", err)
		return nil
	}

	return &ProductService{
		db,
		serviceName,
	}
}

// QueryProducts is a simple SQL WHERE query on a single column
func (s ProductService) QueryProducts(column, term string) ([]spec.Product, error) {
	rows, err := s.Query("SELECT * FROM products WHERE "+column+" = ?", term)
	if err != nil {
		prob := problem.New("err://products-db", "Database query error", 500, err.Error(), s.serviceName)
		return nil, prob
	}

	return s.processRows(rows)
}

// AllProducts returns all products from the DB, yeah this is pretty dumb
func (s ProductService) AllProducts() ([]spec.Product, error) {
	rows, err := s.Query("SELECT * FROM products")
	if err != nil {
		prob := problem.New("err://products-db", "Database query error", 500, err.Error(), s.serviceName)
		return nil, prob
	}

	return s.processRows(rows)
}

// SearchProducts is a text search in name or  product description
func (s ProductService) SearchProducts(query string) ([]spec.Product, error) {
	rows, err := s.Query("SELECT * FROM products WHERE (description LIKE ? OR name LIKE ?)", "%"+query+"%", "%"+query+"%")
	if err != nil {
		prob := problem.New("err://products-db", "Database query error", 500, err.Error(), s.serviceName)
		return nil, prob
	}

	return s.processRows(rows)
}

// Helper function to take a bunch of rows and return as a slice of Products
func (s ProductService) processRows(rows *sql.Rows) ([]spec.Product, error) {
	products := []spec.Product{}
	defer rows.Close()
	for rows.Next() {
		p := spec.Product{}
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Cost, &p.Image, &p.OnOffer)
		if err != nil {
			prob := problem.New("err://products-db", "Error reading row", 500, err.Error(), s.serviceName)
			return nil, prob
		}
		products = append(products, p)
	}

	return products, nil
}

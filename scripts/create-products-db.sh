#!/bin/bash

#
# Build the cmd/products/sqlite.db database file
#

echo "### Using cmd/products/sqlite.db"
echo "### Droping products table"
sqlite3 cmd/products/sqlite.db "DROP TABLE IF EXISTS products"

echo "### Creating products table"
sqlite3 cmd/products/sqlite.db "CREATE TABLE products ( 
  id integer not null primary key,
  name text NOT null,
  description TEXT,
  cost REAL,
  image TEXT,
  onoffer INT);"

echo "### Importing etc/products.csv into products table"
sqlite3 -csv cmd/products/sqlite.db ".import etc/products.csv products"

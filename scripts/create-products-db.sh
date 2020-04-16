#!/bin/bash
sqlite3 cmd/products/sqlite.db "DROP TABLE IF EXISTS products"

sqlite3 cmd/products/sqlite.db "CREATE TABLE products ( 
  id integer not null primary key,
  name text NOT null,
  descripton TEXT,
  cost REAL,
  image TEXT,
  onoffer INT);"

sqlite3 -csv cmd/products/sqlite.db ".import etc/products.csv products"

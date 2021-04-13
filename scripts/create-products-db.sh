#!/bin/bash

#
# Build the cmd/products/sqlite.db database file
#

outputDb=${1:-"cmd/products/sqlite.db"}
inputCsv=${2:-"etc/products.csv"}

echo "🠶🠶🠶 Will create or update: $outputDb"
echo "🠶🠶🠶 Droping products table"
sqlite3 "$outputDb" "DROP TABLE IF EXISTS products"

echo "🠶🠶🠶 Creating products table"
sqlite3 "$outputDb" "CREATE TABLE products ( 
  id TEXT not null primary key,
  name text NOT null,
  description TEXT,
  cost REAL,
  image TEXT,
  onoffer INT);"

echo "🠶🠶🠶 Importing $inputCsv into products table"
sqlite3 -csv "$outputDb" ".import $inputCsv products"

echo "🠶🠶🠶 Database products table contains: $(sqlite3 "$outputDb" 'SELECT COUNT(*) FROM products;') products"
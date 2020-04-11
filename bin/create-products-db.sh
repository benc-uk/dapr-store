#!/bin/bash
cat etc/products.sql | sqlite3 services/products/sqlite.db

drop TABLE products;
create table products (
  id integer not null primary key AUTOINCREMENT,
  name text NOT null,
  descripton TEXT,
  cost REAL,
  image TEXT,
  onoffer INT);
 INSERT INTO products (name, descripton, cost, image, onoffer) 
        VALUES("Hat", "A blue hat", 23.66, "hat.png", 0);
 INSERT INTO products (name, descripton, cost, image, onoffer) 
        VALUES("Lawnmower", "It's green and noisy", 700, "hat.png", 1);
 INSERT INTO products (name, descripton, cost, image, onoffer) 
        VALUES("Cheese", "Yellow and smelly", 3.96, "cheese.png", 0);
 
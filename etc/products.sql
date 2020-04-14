drop TABLE products;
create table products (
  id integer not null primary key AUTOINCREMENT,
  name text NOT null,
  descripton TEXT,
  cost REAL,
  image TEXT,
  onoffer INT);
 INSERT INTO products (name, descripton, cost, image, onoffer) 
        VALUES("Top Hat (6″)", "Made from 100% Wool, Handmade, Inner Black Satin Lining, Inner Leatherette Sweatband", 39.95, "/catalog/1.jpg", 0);
 INSERT INTO products (name, descripton, cost, image, onoffer) 
        VALUES("Paisley Pattern, Silk Bow Tie", "Burgundy, Blue & Silver Paisley Patterned Bow Tie, Ready Tied, Fits neck sizes 28cm to 50cm", 15.00, "/catalog/2.jpg", 1);
 INSERT INTO products (name, descripton, cost, image, onoffer) 
        VALUES("Mens Hipster Paisley Waistcoat", "Paisley pattern, 70% Cotton,30% Polyester", 22.50, "/catalog/3.jpg", 0);
 
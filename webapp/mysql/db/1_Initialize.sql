USE isuumo;

LOAD DATA INFILE '/var/lib/mysql-files/chairData.csv' IGNORE INTO TABLE chair FIELDS TERMINATED BY ',' ENCLOSED BY '"' LINES TERMINATED BY '\n' (id, thumbnail, name, price, height, width, depth, view_count, stock, color, description, features, kind);
LOAD DATA INFILE '/var/lib/mysql-files/estateData.csv' IGNORE INTO TABLE estate FIELDS TERMINATED BY ',' ENCLOSED BY '"' LINES TERMINATED BY '\n' (id, thumbnail, name, latitude, longitude, address, rent, door_height, door_width, view_count, description, features);

UPDATE chair SET view_count = 10000, stock = 2;
UPDATE estate SET view_count = 10000;

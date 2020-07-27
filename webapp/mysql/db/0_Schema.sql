DROP DATABASE IF EXISTS isuumo;
CREATE DATABASE isuumo;

DROP TABLE IF EXISTS isuumo.estate;
DROP TABLE IF EXISTS isuumo.chair;

CREATE TABLE isuumo.estate (
    id INTEGER PRIMARY KEY,
    thumbnail VARCHAR(256) NOT NULL,
    name VARCHAR(64) NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    address VARCHAR(128) NOT NULL,
    rent INTEGER NOT NULL,
    door_height INTEGER NOT NULL,
    door_width INTEGER NOT NULL,
    view_count INTEGER DEFAULT 0 NOT NULL,
    description TEXT NOT NULL,
    features VARCHAR(256) NOT NULL
);

CREATE TABLE isuumo.chair (
    id INTEGER PRIMARY KEY,
    thumbnail TEXT,
    name VARCHAR(64) NOT NULL,
    description TEXT NOT NULL,
    price INTEGER NOT NULL,
    height INTEGER NOT NULL,
    width INTEGER NOT NULL,
    depth INTEGER NOT NULL,
    view_count INTEGER NOT NULL DEFAULT 0,
    stock INTEGER NOT NULL DEFAULT 0,
    color VARCHAR(64) NOT NULL,
    features VARCHAR(64) NOT NULL,
    kind VARCHAR(64) NOT NULL
);

LOAD DATA INFILE '/var/lib/mysql-files/chairData.csv' INTO TABLE isuumo.chair FIELDS TERMINATED BY ',' ENCLOSED BY '"' LINES TERMINATED BY '\n' IGNORE 1 LINES (id, thumbnail, name, price, height, width, depth, view_count, stock, color, description, features, kind);
LOAD DATA INFILE '/var/lib/mysql-files/estateData.csv' INTO TABLE isuumo.estate FIELDS TERMINATED BY ',' ENCLOSED BY '"' LINES TERMINATED BY '\n' IGNORE 1 LINES (id, thumbnail, name, latitude, longitude, address, rent, door_height, door_width, view_count, description, features);

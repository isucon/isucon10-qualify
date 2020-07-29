\encoding UTF8;

DROP DATABASE IF EXISTS isuumo;
CREATE DATABASE isuumo;

\c isuumo;

DROP TABLE IF EXISTS estate;
DROP TABLE IF EXISTS chair;

CREATE TABLE estate (
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

CREATE TABLE chair (
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

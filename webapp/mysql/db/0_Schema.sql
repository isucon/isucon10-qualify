DROP DATABASE IF EXISTS isuumo;
CREATE DATABASE isuumo DEFAULT CHARACTER SET utf8mb4;

USE isuumo;

DROP TABLE IF EXISTS estate;
DROP TABLE IF EXISTS chair;

CREATE TABLE estate (
    `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
    `thumbnail` VARCHAR(256) NOT NULL,
    `name` VARCHAR(64) NOT NULL,
    `latitude` DOUBLE NOT NULL,
    `longitude` DOUBLE NOT NULL,
    `address` VARCHAR(128) NOT NULL,
    `rent` INTEGER NOT NULL,
    `door_height` INTEGER NOT NULL,
    `door_width` INTEGER NOT NULL,
    `view_count` INTEGER DEFAULT 0 NOT NULL,
    `description` TEXT NOT NULL,
    `features` VARCHAR(256) NOT NULL
)ENGINE=InnoDB;

CREATE TABLE chair (
    `id` INTEGER PRIMARY KEY AUTO_INCREMENT,
    `thumbnail` TEXT,
    `name` VARCHAR(64) NOT NULL,
    `description` TEXT NOT NULL,
    `price` INTEGER NOT NULL,
    `height` INTEGER NOT NULL,
    `width` INTEGER NOT NULL,
    `depth` INTEGER NOT NULL,
    `view_count` INTEGER NOT NULL DEFAULT 0,
    `stock` INTEGER NOT NULL DEFAULT 0,
    `color` VARCHAR(64) NOT NULL,
    `features` VARCHAR(64) NOT NULL,
    `kind` VARCHAR(64) NOT NULL
)ENGINE=InnoDB;

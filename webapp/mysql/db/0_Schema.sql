DROP DATABASE IF EXISTS isuumo;
CREATE DATABASE isuumo;

DROP TABLE IF EXISTS isuumo.estate;
DROP TABLE IF EXISTS isuumo.chair;

CREATE TABLE isuumo.estate
(
    id          INTEGER             NOT NULL PRIMARY KEY,
    name        VARCHAR(64)         NOT NULL,
    description VARCHAR(4096)       NOT NULL,
    thumbnail   VARCHAR(128)        NOT NULL,
    address     VARCHAR(128)        NOT NULL,
    latitude    DOUBLE PRECISION    NOT NULL,
    longitude   DOUBLE PRECISION    NOT NULL,
    rent        INTEGER             NOT NULL,
    door_height INTEGER             NOT NULL,
    door_width  INTEGER             NOT NULL,
    longer      INTEGER             GENERATED ALWAYS AS (GREATEST(door_height, door_width)),
    shorter     INTEGER             GENERATED ALWAYS AS (LEAST(door_height, door_width)),
    features    VARCHAR(64)         NOT NULL,
    popularity  INTEGER             NOT NULL,
    rev_pop     INTEGER             GENERATED ALWAYS AS (-1 * popularity)
);

CREATE TABLE isuumo.chair
(
    id          INTEGER         NOT NULL PRIMARY KEY,
    name        VARCHAR(64)     NOT NULL,
    description VARCHAR(4096)   NOT NULL,
    thumbnail   VARCHAR(128)    NOT NULL,
    price       INTEGER         NOT NULL,
    height      INTEGER         NOT NULL,
    width       INTEGER         NOT NULL,
    depth       INTEGER         NOT NULL,
    shortest    INTEGER         GENERATED ALWAYS AS (LEAST(height, width, depth)),
    second      INTEGER         GENERATED ALWAYS AS (IF(height < width, IF(width < depth, width, IF(height < depth, depth, height)), IF(height < depth, height, IF(width < depth, depth, width)))),
    color       VARCHAR(64)     NOT NULL,
    features    VARCHAR(64)     NOT NULL,
    kind        VARCHAR(64)     NOT NULL,
    popularity  INTEGER         NOT NULL,
    rev_pop     INTEGER             GENERATED ALWAYS AS (-1 * popularity),
    stock       INTEGER         NOT NULL
);

CREATE INDEX rent_idx ON isuumo.estate(rent);
CREATE INDEX price_idx ON isuumo.chair(price);
CREATE INDEX popularity_idx ON isuumo.estate(rev_pop ASC);
CREATE INDEX popularity_idx ON isuumo.chair(rev_pop ASC);

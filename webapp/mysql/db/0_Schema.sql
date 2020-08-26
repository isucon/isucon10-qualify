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
    rentRange   INTEGER             GENERATED ALWAYS AS (GREATEST(0, LEAST((rent - 1) / 50000, 3))),
    door_height INTEGER             NOT NULL,
    heightRange INTEGER             GENERATED ALWAYS AS (GREATEST(0, LEAST((door_height - 51) / 30, 3))),
    door_width  INTEGER             NOT NULL,
    widthRange  INTEGER             GENERATED ALWAYS AS (GREATEST(0, LEAST((door_width - 51) / 30, 3))),
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
    priceRange  INTEGER         GENERATED ALWAYS AS (GREATEST(0, LEAST((price - 1) / 3000, 5))),
    height      INTEGER         NOT NULL,
    heightRange INTEGER         GENERATED ALWAYS AS (GREATEST(0, LEAST((height - 51) / 30, 3))),
    width       INTEGER         NOT NULL,
    widthRange  INTEGER         GENERATED ALWAYS AS (GREATEST(0, LEAST((width - 51) / 30, 3))),
    depth       INTEGER         NOT NULL,
    depthRange  INTEGER         GENERATED ALWAYS AS (GREATEST(0, LEAST((depth - 51) / 30, 3))),
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
CREATE INDEX rent_range_idx ON isuumo.estate(rentRange);
CREATE INDEX height_range_idx ON isuumo.estate(heightRange);
CREATE INDEX width_range_idx ON isuumo.estate(widthRange);
CREATE INDEX price_range_idx ON isuumo.chair(priceRange);
CREATE INDEX height_range_idx ON isuumo.chair(heightRange);
CREATE INDEX width_range_idx ON isuumo.chair(widthRange);
CREATE INDEX depth_range_idx ON isuumo.chair(depthRange);
CREATE INDEX color_idx ON isuumo.chair(color);
CREATE INDEX kind_idx ON isuumo.chair(kind);

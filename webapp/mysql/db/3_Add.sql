ALTER TABLE estate ADD pt POINT AS (POINT(latitude,longitude));
ALTER TABLE estate ADD ngpopularity INTEGER AS (-popularity);

ALTER TABLE estate ADD INDEX index_ngpopularity_id (ngpopularity,id);
ALTER TABLE estate ADD INDEX index_door_height_rent (door_height,rent);
ALTER TABLE estate ADD INDEX index_door_width_rent (door_width,rent);
ALTER TABLE estate ADD INDEX index_door_height_door_width (door_height,door_width);
ALTER TABLE estate ADD INDEX index_id (id);
ALTER TABLE estate ADD INDEX index_rent (rent);
ALTER TABLE estate ADD INDEX index_door_height (door_height);
ALTER TABLE estate ADD INDEX index_door_width (door_width);
ALTER TABLE estate ADD INDEX index_pt (pt);

ALTER TABLE chair ADD ngpopularity INTEGER AS (-popularity);
ALTER TABLE chair ADD INDEX index_ngpopularity_id (ngpopularity,id);
ALTER TABLE chair ADD INDEX index_price_stock (price,stock);
ALTER TABLE chair ADD INDEX index_height_stock (height,stock);
ALTER TABLE chair ADD INDEX index_depth_stock (depth,stock);
ALTER TABLE chair ADD INDEX index_features_stock (features,stock);
ALTER TABLE chair ADD INDEX index_kind_stock (kind,stock);
ALTER TABLE chair ADD INDEX index_id (id);
ALTER TABLE chair ADD INDEX index_price (price);
ALTER TABLE chair ADD INDEX index_height (height);
ALTER TABLE chair ADD INDEX index_width (width);
ALTER TABLE chair ADD INDEX index_depth (depth);
ALTER TABLE chair ADD INDEX index_color (color);
ALTER TABLE chair ADD INDEX index_features (features);
ALTER TABLE chair ADD INDEX index_kind (kind);
ALTER TABLE chair ADD INDEX index_stock (stock);
ALTER TABLE chair ADD INDEX chair_price_id_idx(price, id);
ALTER TABLE estate ADD INDEX estate_rent_id_idx(rent, id);
ALTER TABLE chair ADD INDEX chair_popularity_id_idx(popularity_desc, id);
ALTER TABLE estate ADD INDEX estate_popularity_id_idx(popularity_desc, id);
ALTER TABLE estate ADD SPATIAL INDEX estate_point_idx(point);

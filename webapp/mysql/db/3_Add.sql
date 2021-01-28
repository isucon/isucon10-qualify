ALTER TABLE estate ADD pt POINT AS (POINT(latitude,longitude));
ALTER TABLE estate ADD ngpopularity INTEGER AS (-popularity);
ALTER TABLE estate ADD INDEX ngpopularity_id (ngpopularity,id);
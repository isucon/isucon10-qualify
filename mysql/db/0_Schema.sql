use isuumo;

create table estate (
    id int auto_increment not null,
    thumbnails varchar(256) not null,
    name varchar(64) not null,
    latitude double not null,
    longitude double not null,
    address varchar(128) not null,
    rent integer not null,
    door_height integer not null,
    door_width integer not null,
    view_count integer default 0 not null,
    description text not null,
    features varchar(256) not null,
    primary key(id)
)ENGINE=InnoDB;

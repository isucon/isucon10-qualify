use isuumo;

create table estate (
    id integer primary key auto_increment,
    thumbnail varchar(256) not null,
    name varchar(64) not null,
    latitude double not null,
    longitude double not null,
    address varchar(128) not null,
    rent integer not null,
    door_height integer not null,
    door_width integer not null,
    view_count integer default 0 not null,
    description text not null,
    features varchar(256) not null
)ENGINE=InnoDB;

create table chair (
    id integer primary key auto_increment,
    thumbnail text,
    name varchar(64) not null,
    description text not null,
    price integer not null,
    height integer not null,
    width integer not null,
    depth integer not null,
    view_count integer not null default 0,
    stock integer not null default 0,
    color varchar(64) not null,
    features varchar(64) not null,
    kind varchar(64) not null
)ENGINE=InnoDB;

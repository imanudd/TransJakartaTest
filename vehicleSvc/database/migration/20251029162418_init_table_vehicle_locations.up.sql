-- +migrate Up
CREATE TABLE vehicle_locations(
    vehicle_id varchar(100) not null,
    latitude double precision not null,
    longitude double precision not null,
    timestamp int not null
);
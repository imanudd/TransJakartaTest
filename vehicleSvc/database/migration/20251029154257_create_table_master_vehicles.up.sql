-- +migrate Up
CREATE TABLE vehicles(
                         vehicle_id varchar(100) not null,
                         vehicle_type varchar(100) not null,
                         vehicle_code varchar(100) not null
);
create table digimap_db.indoormap_api_event_locations
(
    id          bigserial
        constraint idx_21344_primary
            primary key,
    event_id    char(32) not null
        constraint indoormap_api_event__event_id_35d8fa44_fk_indoormap
            references digimap_db.indoormap_api_event,
    location_id char(32) not null
        constraint indoormap_api_event__location_id_9e57c39f_fk_indoormap
            references digimap_db.indoormap_api_location
);

alter table digimap_db.indoormap_api_event_locations
    owner to postgres;

create index idx_21344_indoormap_api_event__location_id_9e57c39f_fk_indoorma
    on digimap_db.indoormap_api_event_locations (location_id);

create unique index idx_21344_indoormap_api_event_locations_event_id_location_id_54
    on digimap_db.indoormap_api_event_locations (event_id, location_id);


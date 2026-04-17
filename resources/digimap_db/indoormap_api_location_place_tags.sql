create table digimap_db.indoormap_api_location_place_tags
(
    id          bigserial
        constraint idx_21460_primary
            primary key,
    location_id char(32) not null
        constraint indoormap_api_locati_location_id_b40aef29_fk_indoormap
            references digimap_db.indoormap_api_location,
    tag_id      char(32) not null
        constraint indoormap_api_locati_tag_id_9e781470_fk_indoormap
            references digimap_db.indoormap_api_tag
);

alter table digimap_db.indoormap_api_location_place_tags
    owner to postgres;

create index idx_21460_indoormap_api_locati_tag_id_9e781470_fk_indoormap
    on digimap_db.indoormap_api_location_place_tags (tag_id);

create unique index idx_21460_indoormap_api_location_tags_location_id_tag_id_bb6cda
    on digimap_db.indoormap_api_location_place_tags (location_id, tag_id);


create table digimap_db.indoormap_api_location_common_categories
(
    id                  bigserial
        constraint idx_21452_primary
            primary key,
    location_id         char(32) not null
        constraint indoormap_api_locati_location_id_5431b8ac_fk_indoormap
            references digimap_db.indoormap_api_location,
    locationcategory_id char(32) not null
        constraint indoormap_api_locati_locationcategory_id_b3cb6515_fk_indoormap
            references digimap_db.indoormap_api_locationcategory
);

alter table digimap_db.indoormap_api_location_common_categories
    owner to postgres;

create index idx_21452_indoormap_api_locati_locationcategory_id_b3cb6515_fk_
    on digimap_db.indoormap_api_location_common_categories (locationcategory_id);

create unique index idx_21452_indoormap_api_location_c_location_id_locationcate_d82
    on digimap_db.indoormap_api_location_common_categories (location_id, locationcategory_id);


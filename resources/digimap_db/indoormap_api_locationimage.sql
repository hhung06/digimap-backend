create table digimap_db.indoormap_api_locationimage
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32) not null
        constraint idx_21484_primary
            primary key,
    location_id    char(32)
        constraint indoormap_api_locati_location_id_74cf1f87_fk_indoormap
            references digimap_db.indoormap_api_location,
    large          varchar(1000),
    medium         varchar(1000),
    original       varchar(1000),
    small          varchar(1000),
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_locationimage
    owner to postgres;

create index idx_21484_indoormap_api_locati_location_id_74cf1f87_fk_indoorma
    on digimap_db.indoormap_api_locationimage (location_id);


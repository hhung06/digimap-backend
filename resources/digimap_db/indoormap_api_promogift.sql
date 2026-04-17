create table digimap_db.indoormap_api_promogift
(
    deleted_at     timestamp with time zone,
    restored_at    timestamp with time zone,
    transaction_id char(32),
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    id             char(32)     not null
        constraint idx_21654_primary
            primary key,
    gift_name      varchar(255) not null,
    gift_details   text,
    venue_id       char(32)
        constraint indoormap_api_promog_venue_id_00e7e69d_fk_indoormap
            references digimap_db.indoormap_api_venue,
    external_id    varchar(255),
    image          varchar(1000),
    localization   json,
    gift_counts    bigint,
    gift_quantity  bigint,
    location_id    char(32)
        constraint indoormap_api_promog_location_id_84598454_fk_indoormap
            references digimap_db.indoormap_api_location
);

alter table digimap_db.indoormap_api_promogift
    owner to postgres;

create index idx_21654_indoormap_api_promog_venue_id_00e7e69d_fk_indoormap
    on digimap_db.indoormap_api_promogift (venue_id);

create index idx_21654_indoormap_api_promog_location_id_84598454_fk_indoorma
    on digimap_db.indoormap_api_promogift (location_id);


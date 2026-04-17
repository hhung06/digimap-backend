create table digimap_db.indoormap_api_product
(
    created_at       timestamp with time zone,
    updated_at       timestamp with time zone,
    deleted_at       timestamp with time zone,
    id               char(32)   not null
        constraint idx_21578_primary
            primary key,
    name             text,
    code             text,
    size             text,
    price            text,
    origin_country   text,
    expiration       text,
    description      text,
    custom           json,
    location_id      char(32)
        constraint indoormap_api_produc_location_id_edf95cb2_fk_indoormap
            references digimap_db.indoormap_api_location,
    main_category_id char(32)
        constraint indoormap_api_produc_main_category_id_7221dfe3_fk_indoormap
            references digimap_db.indoormap_api_productcategory,
    localization     json,
    source           varchar(8) not null,
    image            varchar(255),
    restored_at      timestamp with time zone,
    transaction_id   char(32),
    venue_id         char(32)
        constraint indoormap_api_produc_venue_id_927295fe_fk_indoormap
            references digimap_db.indoormap_api_venue
);

alter table digimap_db.indoormap_api_product
    owner to postgres;

create index idx_21578_indoormap_api_produc_main_category_id_7221dfe3_fk_ind
    on digimap_db.indoormap_api_product (main_category_id);

create index idx_21578_indoormap_api_produc_location_id_edf95cb2_fk_indoorma
    on digimap_db.indoormap_api_product (location_id);

create index idx_21578_indoormap_api_produc_venue_id_927295fe_fk_indoormap
    on digimap_db.indoormap_api_product (venue_id);


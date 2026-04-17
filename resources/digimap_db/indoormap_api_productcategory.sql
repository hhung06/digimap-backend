create table digimap_db.indoormap_api_productcategory
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32)      not null
        constraint idx_21600_primary
            primary key,
    name           varchar(1000) not null,
    venue_id       char(32)
        constraint indoormap_api_produc_venue_id_348a595a_fk_indoormap
            references digimap_db.indoormap_api_venue,
    localization   json,
    source         varchar(8)    not null,
    externalid     varchar(255),
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_productcategory
    owner to postgres;

create index idx_21600_indoormap_api_produc_venue_id_348a595a_fk_indoormap
    on digimap_db.indoormap_api_productcategory (venue_id);


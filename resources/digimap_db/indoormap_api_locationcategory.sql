create table digimap_db.indoormap_api_locationcategory
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32)   not null
        constraint idx_21471_primary
            primary key,
    name           varchar(100),
    icon           varchar(1000),
    venue_id       char(32)
        constraint indoormap_api_locati_venue_id_6a5feeb8_fk_indoormap
            references digimap_db.indoormap_api_venue,
    color          varchar(255),
    sortindex      bigint     not null,
    icondefault    varchar(1000),
    visible        boolean    not null,
    localization   json,
    source         varchar(8) not null,
    type           varchar(255),
    shortname      varchar(255),
    restored_at    timestamp with time zone,
    transaction_id char(32),
    externalid     varchar(255),
    description    text,
    image          varchar(1000)
);

alter table digimap_db.indoormap_api_locationcategory
    owner to postgres;

create index idx_21471_indoormap_api_locati_venue_id_6a5feeb8_fk_indoormap
    on digimap_db.indoormap_api_locationcategory (venue_id);


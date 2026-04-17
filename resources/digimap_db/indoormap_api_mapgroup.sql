create table digimap_db.indoormap_api_mapgroup
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32) not null
        constraint idx_21539_primary
            primary key,
    type           varchar(255),
    name           varchar(100),
    shortname      varchar(50),
    sortindex      bigint   not null,
    venue_id       char(32)
        constraint indoormap_api_mapgro_venue_id_ae910a39_fk_indoormap
            references digimap_db.indoormap_api_venue,
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_mapgroup
    owner to postgres;

create index idx_21539_indoormap_api_mapgro_venue_id_ae910a39_fk_indoormap
    on digimap_db.indoormap_api_mapgroup (venue_id);


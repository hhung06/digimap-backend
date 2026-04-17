create table digimap_db.indoormap_api_connectionlevel
(
    id             bigserial
        constraint idx_21310_primary
            primary key,
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    active         boolean not null,
    connection_id  char(32)
        constraint indoormap_api_connec_connection_id_4f11bb58_fk_indoormap
            references digimap_db.indoormap_api_connection,
    level_id       char(32)
        constraint indoormap_api_connec_level_id_91d5e466_fk_indoormap
            references digimap_db.indoormap_api_level,
    element_id     char(32),
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_connectionlevel
    owner to postgres;

create index idx_21310_indoormap_api_connec_level_id_91d5e466_fk_indoormap
    on digimap_db.indoormap_api_connectionlevel (level_id);

create index idx_21310_indoormap_api_connec_connection_id_4f11bb58_fk_indoor
    on digimap_db.indoormap_api_connectionlevel (connection_id);


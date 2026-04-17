create table digimap_db.indoormap_api_georeference
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32)         not null
        constraint idx_21393_primary
            primary key,
    controlx       integer          not null,
    controly       integer          not null,
    targetx        double precision not null,
    targety        double precision not null,
    level_id       char(32)
        constraint indoormap_api_georef_level_id_d043778c_fk_indoormap
            references digimap_db.indoormap_api_level,
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_georeference
    owner to postgres;

create index idx_21393_indoormap_api_georef_level_id_d043778c_fk_indoormap
    on digimap_db.indoormap_api_georeference (level_id);


create table digimap_db.indoormap_api_syncdata
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32)   not null
        constraint idx_21737_primary
            primary key,
    venue_id       char(32),
    level_id       char(32),
    version        integer    not null,
    state          varchar(6) not null,
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_syncdata
    owner to postgres;


create table digimap_db.indoormap_api_asset
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             varchar(255) not null
        constraint idx_21272_primary
            primary key,
    file_name      varchar(1000),
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_asset
    owner to postgres;


create table digimap_db.indoormap_api_leveltype
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32) not null
        constraint idx_21424_primary
            primary key,
    name           varchar(20),
    value          integer  not null,
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_leveltype
    owner to postgres;


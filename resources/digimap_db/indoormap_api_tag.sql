create table digimap_db.indoormap_api_tag
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32) not null
        constraint idx_21743_primary
            primary key,
    name           varchar(100),
    localization   json,
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_tag
    owner to postgres;


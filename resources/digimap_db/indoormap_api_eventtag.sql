create table digimap_db.indoormap_api_eventtag
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32) not null
        constraint idx_21366_primary
            primary key,
    name           varchar(255),
    localization   json,
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_eventtag
    owner to postgres;


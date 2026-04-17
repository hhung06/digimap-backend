create table digimap_db.indoormap_api_libraryasset
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32)    not null
        constraint idx_21429_primary
            primary key,
    type           varchar(11) not null,
    file_name      varchar(1000),
    thumbnail      varchar(1000),
    status         varchar(11) not null,
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_libraryasset
    owner to postgres;


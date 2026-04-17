create table digimap_db.indoormap_api_new
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32) not null
        constraint idx_21546_primary
            primary key,
    title          varchar(255),
    bannerimage    varchar(100),
    contenturl     varchar(255),
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_new
    owner to postgres;


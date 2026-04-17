create table digimap_db.indoormap_api_customer
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32) not null
        constraint idx_21331_primary
            primary key,
    name           varchar(100),
    phone          varchar(20),
    url            varchar(100),
    description    text,
    restored_at    timestamp with time zone,
    transaction_id char(32),
    address        varchar(100),
    email          varchar(255),
    image          varchar(100)
);

alter table digimap_db.indoormap_api_customer
    owner to postgres;


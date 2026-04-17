create table digimap_db.indoormap_api_connection
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32)         not null
        constraint idx_21296_primary
            primary key,
    name           varchar(255),
    type           smallint         not null,
    accessible     boolean          not null,
    venue_id       char(32)
        constraint indoormap_api_connec_venue_id_3002b720_fk_indoormap
            references digimap_db.indoormap_api_venue,
    state          smallint         not null,
    status         smallint         not null,
    x              double precision not null,
    y              double precision not null,
    active         boolean          not null,
    externalid     varchar(255),
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_connection
    owner to postgres;

create index idx_21296_indoormap_api_connec_venue_id_3002b720_fk_indoormap
    on digimap_db.indoormap_api_connection (venue_id);


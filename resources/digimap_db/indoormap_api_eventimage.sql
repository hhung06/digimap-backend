create table digimap_db.indoormap_api_eventimage
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32) not null
        constraint idx_21359_primary
            primary key,
    image          varchar(1000),
    event_id       char(32) not null
        constraint indoormap_api_eventi_event_id_7c1f1f54_fk_indoormap
            references digimap_db.indoormap_api_event,
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_eventimage
    owner to postgres;

create index idx_21359_indoormap_api_eventi_event_id_7c1f1f54_fk_indoormap
    on digimap_db.indoormap_api_eventimage (event_id);


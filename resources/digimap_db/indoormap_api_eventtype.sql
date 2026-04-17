create table digimap_db.indoormap_api_eventtype
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32) not null
        constraint idx_21372_primary
            primary key,
    name           varchar(255),
    venue_id       char(32)
        constraint indoormap_api_eventt_venue_id_aa1577e0_fk_indoormap
            references digimap_db.indoormap_api_venue,
    localization   json,
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_eventtype
    owner to postgres;

create index idx_21372_indoormap_api_eventt_venue_id_aa1577e0_fk_indoormap
    on digimap_db.indoormap_api_eventtype (venue_id);


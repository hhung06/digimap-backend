create table digimap_db.indoormap_api_event
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32) not null
        constraint idx_21337_primary
            primary key,
    bannerimage    varchar(1000),
    title          varchar(255),
    starttime      timestamp with time zone,
    endtime        timestamp with time zone,
    type_id        char(32),
    contentdetail  text,
    contenturl     varchar(255),
    description    text,
    iconimage      varchar(1000),
    showendtime    timestamp with time zone,
    showstarttime  timestamp with time zone,
    venue_id       char(32)
        constraint indoormap_api_event_venue_id_08a7d602_fk_indoormap_api_venue_id
            references digimap_db.indoormap_api_venue,
    localization   json,
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_event
    owner to postgres;

create index idx_21337_indoormap_api_event_type_id_19aeda3f
    on digimap_db.indoormap_api_event (type_id);

create index idx_21337_indoormap_api_event_venue_id_08a7d602_fk_indoormap_ap
    on digimap_db.indoormap_api_event (venue_id);


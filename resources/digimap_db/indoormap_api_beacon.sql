create table digimap_db.indoormap_api_beacon
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32)         not null
        constraint idx_21278_primary
            primary key,
    hwid           varchar(100),
    vendorkey      varchar(100),
    lotkey         varchar(100),
    uuid           varchar(100),
    mac            varchar(100),
    radius         integer          not null,
    battery        integer          not null,
    name           varchar(100),
    positionx      double precision not null,
    positiony      double precision not null,
    level          integer          not null,
    isenable       boolean          not null,
    major          integer,
    minor          integer,
    voltage        integer,
    txpower        integer,
    venue_id       char(32)
        constraint indoormap_api_beacon_venue_id_08aabc2f_fk_indoormap_api_venue_i
            references digimap_db.indoormap_api_venue,
    level_id       char(32)
        constraint indoormap_api_beacon_level_id_029d21be_fk_indoormap_api_level_i
            references digimap_db.indoormap_api_level,
    element_id     char(32),
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_beacon
    owner to postgres;

create index idx_21278_indoormap_api_beacon_level_id_029d21be_fk_indoormap_a
    on digimap_db.indoormap_api_beacon (level_id);

create index idx_21278_indoormap_api_beacon_venue_id_08aabc2f_fk_indoormap_a
    on digimap_db.indoormap_api_beacon (venue_id);


create table digimap_db.indoormap_api_level
(
    id             char(32)         not null
        constraint idx_21409_primary
            primary key,
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    externalid     varchar(255),
    height         integer          not null,
    scale          double precision not null,
    shortname      varchar(255),
    venue_id       char(32)
        constraint indoormap_api_level_venue_id_2413c37c_fk
            references digimap_db.indoormap_api_venue,
    width          integer          not null,
    perspective_id char(32)
        constraint indoormap_api_level_perspective_id_2e9d9229_fk_indoormap
            references digimap_db.indoormap_api_perspective,
    type_id        char(32)
        constraint indoormap_api_level_type_id_a18364cf_fk_indoormap
            references digimap_db.indoormap_api_leveltype,
    name           varchar(255),
    publish        boolean          not null,
    bearing        double precision not null,
    latitude       double precision not null,
    longitude      double precision not null,
    level_height   double precision not null,
    level_width    double precision not null,
    file_ids       text,
    elevation      integer,
    mapgroup_id    char(32)
        constraint indoormap_api_level_mapgroup_id_9d651189_fk_indoormap
            references digimap_db.indoormap_api_mapgroup,
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_level
    owner to postgres;

create unique index idx_21409_perspective_id
    on digimap_db.indoormap_api_level (perspective_id);

create index idx_21409_indoormap_api_level_type_id_a18364cf
    on digimap_db.indoormap_api_level (type_id);

create index idx_21409_indoormap_api_level_mapgroup_id_9d651189_fk_indoormap
    on digimap_db.indoormap_api_level (mapgroup_id);

create index idx_21409_indoormap_api_level_venue_id_2413c37c_fk
    on digimap_db.indoormap_api_level (venue_id);


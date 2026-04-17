create table digimap_db.indoormap_api_venueuserrole
(
    id             bigserial
        constraint idx_21816_primary
            primary key,
    deleted_at     timestamp with time zone,
    restored_at    timestamp with time zone,
    transaction_id char(32),
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    role_type      varchar(20) not null,
    profile_id     integer     not null
        constraint indoormap_api_venueu_profile_id_37cc1c83_fk_indoormap
            references digimap_db.indoormap_api_profile,
    venue_id       char(32)    not null
        constraint indoormap_api_venueu_venue_id_3ab36cdc_fk_indoormap
            references digimap_db.indoormap_api_venue
);

alter table digimap_db.indoormap_api_venueuserrole
    owner to postgres;

create index idx_21816_indoormap_api_venueuserrole_profile_id_37cc1c83
    on digimap_db.indoormap_api_venueuserrole (profile_id);

create index idx_21816_indoormap_api_venueuserrole_venue_id_3ab36cdc
    on digimap_db.indoormap_api_venueuserrole (venue_id);


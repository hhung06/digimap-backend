create table digimap_db.indoormap_api_venueinvitation
(
    deleted_at         timestamp with time zone,
    restored_at        timestamp with time zone,
    transaction_id     char(32),
    created_at         timestamp with time zone,
    updated_at         timestamp with time zone,
    id                 char(32)     not null
        constraint idx_21794_primary
            primary key,
    invited_email      varchar(254) not null,
    role               varchar(20)  not null,
    status             varchar(20)  not null,
    expires_at         timestamp with time zone,
    accepted_at        timestamp with time zone,
    cancelled_at       timestamp with time zone,
    invited_by_id      integer      not null
        constraint indoormap_api_venuei_invited_by_id_2915f3a7_fk_indoormap
            references digimap_db.indoormap_api_profile,
    invited_profile_id integer
        constraint indoormap_api_venuei_invited_profile_id_203f3f41_fk_indoormap
            references digimap_db.indoormap_api_profile,
    venue_id           char(32)     not null
        constraint indoormap_api_venuei_venue_id_4fb38cff_fk_indoormap
            references digimap_db.indoormap_api_venue
);

alter table digimap_db.indoormap_api_venueinvitation
    owner to postgres;

create index idx_21794_indoormap_api_venueinvitation_invited_profile_id_203f
    on digimap_db.indoormap_api_venueinvitation (invited_profile_id);

create index idx_21794_indoormap_api_venueinvitation_venue_id_4fb38cff
    on digimap_db.indoormap_api_venueinvitation (venue_id);

create index idx_21794_indoormap_api_venuei_invited_by_id_2915f3a7_fk_indoor
    on digimap_db.indoormap_api_venueinvitation (invited_by_id);


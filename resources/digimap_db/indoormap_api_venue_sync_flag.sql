create table digimap_db.indoormap_api_venue_sync_flag
(
    deleted_at          timestamp with time zone,
    restored_at         timestamp with time zone,
    transaction_id      char(32),
    created_at          timestamp with time zone,
    updated_at          timestamp with time zone,
    id                  char(32) not null
        constraint idx_21788_primary
            primary key,
    venue_id            char(32) not null,
    has_pending_changes boolean  not null,
    change_source       varchar(20),
    synced_at           timestamp with time zone,
    expo_id             varchar(255)
);

alter table digimap_db.indoormap_api_venue_sync_flag
    owner to postgres;

create index idx_21788_indoormap_a_venue_i_9413cf_idx
    on digimap_db.indoormap_api_venue_sync_flag (venue_id, created_at);

create index idx_21788_indoormap_a_venue_i_c33d5c_idx
    on digimap_db.indoormap_api_venue_sync_flag (venue_id, has_pending_changes);

create index idx_21788_indoormap_api_venue_sync_flag_venue_id_821c02fb
    on digimap_db.indoormap_api_venue_sync_flag (venue_id);


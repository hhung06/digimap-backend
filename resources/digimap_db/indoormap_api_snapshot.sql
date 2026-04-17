create table digimap_db.indoormap_api_snapshot
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32)   not null
        constraint idx_21706_primary
            primary key,
    state          varchar(6) not null,
    venue_id       char(32),
    publish_at     timestamp with time zone,
    publish_by_id  integer
        constraint indoormap_api_snapshot_publish_by_id_4e82db52_fk_auth_user_id
            references digimap_db.auth_user,
    method         integer    not null,
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_snapshot
    owner to postgres;

create index idx_21706_indoormap_api_snapshot_publish_by_id_4e82db52_fk_auth
    on digimap_db.indoormap_api_snapshot (publish_by_id);


create table digimap_db.indoormap_api_appuserdevice
(
    id             bigserial
        constraint idx_21236_primary
            primary key,
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    app            varchar(100),
    device_id      varchar(255),
    fcm_token      text,
    platform       varchar(20),
    is_active      boolean                  not null,
    last_login     timestamp with time zone not null,
    user_id        bigint                   not null
        constraint indoormap_api_appuse_user_id_99c19ad8_fk_indoormap
            references digimap_db.indoormap_api_appuser,
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_appuserdevice
    owner to postgres;

create index idx_21236_indoormap_a_app_387426_idx
    on digimap_db.indoormap_api_appuserdevice (app, is_active);

create index idx_21236_indoormap_a_user_id_2712c9_idx
    on digimap_db.indoormap_api_appuserdevice (user_id, is_active);

create index idx_21236_indoormap_a_last_lo_06a6fc_idx
    on digimap_db.indoormap_api_appuserdevice (last_login);


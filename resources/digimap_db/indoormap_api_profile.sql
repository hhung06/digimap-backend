create table digimap_db.indoormap_api_profile
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    user_id        integer not null
        constraint idx_21647_primary
            primary key
        constraint indoormap_api_profile_user_id_5a841cea_fk_auth_user_id
            references digimap_db.auth_user,
    full_name      varchar(100),
    profile_image  varchar(1000),
    restored_at    timestamp with time zone,
    transaction_id char(32),
    customer_id    char(32)
        constraint indoormap_api_profil_customer_id_f1024b17_fk_indoormap
            references digimap_db.indoormap_api_customer,
    is_partner     boolean not null
);

alter table digimap_db.indoormap_api_profile
    owner to postgres;

create index idx_21647_indoormap_api_profil_customer_id_f1024b17_fk_indoorma
    on digimap_db.indoormap_api_profile (customer_id);


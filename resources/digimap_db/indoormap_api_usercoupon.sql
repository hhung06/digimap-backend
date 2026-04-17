create table digimap_db.indoormap_api_usercoupon
(
    deleted_at     timestamp with time zone,
    restored_at    timestamp with time zone,
    transaction_id char(32),
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    id             char(32)     not null
        constraint idx_21763_primary
            primary key,
    "user"         varchar(255) not null,
    device_id      varchar(255) not null,
    is_used        boolean      not null,
    used_at        timestamp with time zone,
    coupon_id      char(32)     not null
        constraint indoormap_api_userco_coupon_id_a03600a8_fk_indoormap
            references digimap_db.indoormap_api_coupon
);

alter table digimap_db.indoormap_api_usercoupon
    owner to postgres;

create index idx_21763_indoormap_api_userco_coupon_id_a03600a8_fk_indoormap
    on digimap_db.indoormap_api_usercoupon (coupon_id);


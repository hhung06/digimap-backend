create table digimap_db.indoormap_api_coupon_gifts
(
    id           bigserial
        constraint idx_21324_primary
            primary key,
    coupon_id    char(32) not null
        constraint indoormap_api_coupon_coupon_id_b38fa942_fk_indoormap
            references digimap_db.indoormap_api_coupon,
    promogift_id char(32) not null
        constraint indoormap_api_coupon_promogift_id_e7f1bbb9_fk_indoormap
            references digimap_db.indoormap_api_promogift
);

alter table digimap_db.indoormap_api_coupon_gifts
    owner to postgres;

create index idx_21324_indoormap_api_coupon_promogift_id_e7f1bbb9_fk_indoorm
    on digimap_db.indoormap_api_coupon_gifts (promogift_id);

create unique index idx_21324_indoormap_api_coupon_gifts_coupon_id_promogift_id_f82
    on digimap_db.indoormap_api_coupon_gifts (coupon_id, promogift_id);


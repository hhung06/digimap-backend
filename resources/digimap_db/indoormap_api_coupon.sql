create table digimap_db.indoormap_api_coupon
(
    deleted_at     timestamp with time zone,
    restored_at    timestamp with time zone,
    transaction_id char(32),
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    id             char(32)    not null
        constraint idx_21316_primary
            primary key,
    coupon_code    varchar(100),
    status         varchar(20) not null,
    issued_at      timestamp with time zone,
    expired_at     timestamp with time zone,
    survey_id      varchar(255)
        constraint indoormap_api_coupon_survey_id_5c6684cd_fk_indoormap
            references digimap_db.indoormap_api_survey,
    venue_id       char(32)
        constraint indoormap_api_coupon_venue_id_4e51a9b3_fk_indoormap_api_venue_i
            references digimap_db.indoormap_api_venue,
    coupon_name    varchar(255),
    external_id    varchar(255),
    localization   json
);

alter table digimap_db.indoormap_api_coupon
    owner to postgres;

create index idx_21316_indoormap_api_coupon_venue_id_4e51a9b3_fk_indoormap_a
    on digimap_db.indoormap_api_coupon (venue_id);

create index idx_21316_indoormap_api_coupon_survey_id_5c6684cd_fk_indoormap
    on digimap_db.indoormap_api_coupon (survey_id);


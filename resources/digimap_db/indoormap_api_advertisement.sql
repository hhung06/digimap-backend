create table digimap_db.indoormap_api_advertisement
(
    created_at        timestamp with time zone,
    updated_at        timestamp with time zone,
    deleted_at        timestamp with time zone,
    id                char(32)     not null
        constraint idx_21215_primary
            primary key,
    content_cta_url   varchar(255),
    type              varchar(50)  not null,
    end_at            timestamp with time zone,
    venue_id          char(32)
        constraint indoormap_api_advert_venue_id_20557289_fk_indoormap
            references digimap_db.indoormap_api_venue,
    content_image_url varchar(1000),
    displayduration   bigint,
    location_id       char(32)
        constraint indoormap_api_advert_location_id_c8af4309_fk_indoormap
            references digimap_db.indoormap_api_location,
    placement         varchar(255) not null,
    reward_amount     bigint,
    reward_type       varchar(255),
    size_height       bigint,
    size_width        bigint,
    restored_at       timestamp with time zone,
    transaction_id    char(32),
    published_at      timestamp with time zone,
    start_at          timestamp with time zone,
    status            varchar(50)  not null,
    navigate          varchar(50),
    article_id        char(32)
        constraint indoormap_api_advert_article_id_5edbc569_fk_indoormap
            references digimap_db.indoormap_api_article
);

alter table digimap_db.indoormap_api_advertisement
    owner to postgres;

create index idx_21215_indoormap_a_created_4663bd_idx
    on digimap_db.indoormap_api_advertisement (created_at);

create index idx_21215_indoormap_a_venue_i_a00c1c_idx
    on digimap_db.indoormap_api_advertisement (venue_id);

create index idx_21215_indoormap_api_advert_venue_id_20557289_fk_indoormap
    on digimap_db.indoormap_api_advertisement (venue_id);

create index idx_21215_indoormap_a_locatio_7f7354_idx
    on digimap_db.indoormap_api_advertisement (location_id);

create index idx_21215_indoormap_a_type_9c1ce9_idx
    on digimap_db.indoormap_api_advertisement (type);

create index idx_21215_indoormap_api_advert_article_id_5edbc569_fk_indoormap
    on digimap_db.indoormap_api_advertisement (article_id);


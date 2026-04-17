create table digimap_db.indoormap_api_categoryad
(
    deleted_at     timestamp with time zone,
    restored_at    timestamp with time zone,
    transaction_id char(32),
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    id             char(32) not null
        constraint idx_21290_primary
            primary key,
    introduction   varchar(255),
    image          varchar(1000),
    navigate       varchar(50),
    localization   json,
    category_id    char(32)
        constraint indoormap_api_catego_category_id_54655924_fk_indoormap
            references digimap_db.indoormap_api_locationcategory,
    location_id    char(32)
        constraint indoormap_api_catego_location_id_7d4d0c69_fk_indoormap
            references digimap_db.indoormap_api_location,
    venue_id       char(32)
        constraint indoormap_api_catego_venue_id_6566f9e2_fk_indoormap
            references digimap_db.indoormap_api_venue,
    article_id     char(32)
        constraint indoormap_api_catego_article_id_e99c64c9_fk_indoormap
            references digimap_db.indoormap_api_article
);

alter table digimap_db.indoormap_api_categoryad
    owner to postgres;

create index idx_21290_indoormap_api_catego_article_id_e99c64c9_fk_indoormap
    on digimap_db.indoormap_api_categoryad (article_id);

create index idx_21290_indoormap_api_catego_location_id_7d4d0c69_fk_indoorma
    on digimap_db.indoormap_api_categoryad (location_id);

create index idx_21290_indoormap_api_catego_venue_id_6566f9e2_fk_indoormap
    on digimap_db.indoormap_api_categoryad (venue_id);

create index idx_21290_indoormap_api_catego_category_id_54655924_fk_indoorma
    on digimap_db.indoormap_api_categoryad (category_id);


create table digimap_db.indoormap_api_articleimage
(
    deleted_at     timestamp with time zone,
    restored_at    timestamp with time zone,
    transaction_id char(32),
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    id             char(32)      not null
        constraint idx_21256_primary
            primary key,
    image          varchar(1000) not null,
    "order"        bigint        not null,
    article_id     char(32)      not null
        constraint indoormap_api_articl_article_id_2c093e47_fk_indoormap
            references digimap_db.indoormap_api_article
);

alter table digimap_db.indoormap_api_articleimage
    owner to postgres;

create index idx_21256_indoormap_api_articl_article_id_2c093e47_fk_indoormap
    on digimap_db.indoormap_api_articleimage (article_id);


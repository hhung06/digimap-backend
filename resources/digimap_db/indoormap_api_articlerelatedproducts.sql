create table digimap_db.indoormap_api_articlerelatedproducts
(
    deleted_at     timestamp with time zone,
    restored_at    timestamp with time zone,
    transaction_id char(32),
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    id             char(32) not null
        constraint idx_21265_primary
            primary key,
    "order"        bigint   not null,
    article_id     char(32) not null
        constraint indoormap_api_articl_article_id_194a6401_fk_indoormap
            references digimap_db.indoormap_api_article,
    product_id     char(32) not null
        constraint indoormap_api_articl_product_id_c3a31e86_fk_indoormap
            references digimap_db.indoormap_api_product
);

alter table digimap_db.indoormap_api_articlerelatedproducts
    owner to postgres;

create index idx_21265_indoormap_api_articl_product_id_c3a31e86_fk_indoormap
    on digimap_db.indoormap_api_articlerelatedproducts (product_id);

create index idx_21265_indoormap_api_articl_article_id_194a6401_fk_indoormap
    on digimap_db.indoormap_api_articlerelatedproducts (article_id);


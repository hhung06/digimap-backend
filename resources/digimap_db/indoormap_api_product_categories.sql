create table digimap_db.indoormap_api_product_categories
(
    id                 bigserial
        constraint idx_21586_primary
            primary key,
    product_id         char(32) not null
        constraint indoormap_api_produc_product_id_5c419598_fk_indoormap
            references digimap_db.indoormap_api_product,
    productcategory_id char(32) not null
        constraint indoormap_api_produc_productcategory_id_678f4361_fk_indoormap
            references digimap_db.indoormap_api_productcategory
);

alter table digimap_db.indoormap_api_product_categories
    owner to postgres;

create unique index idx_21586_indoormap_api_product_ca_product_id_productcatego_baf
    on digimap_db.indoormap_api_product_categories (product_id, productcategory_id);

create index idx_21586_indoormap_api_produc_productcategory_id_678f4361_fk_i
    on digimap_db.indoormap_api_product_categories (productcategory_id);


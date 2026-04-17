create table digimap_db.indoormap_api_productattachment
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32)    not null
        constraint idx_21593_primary
            primary key,
    title          varchar(255),
    file_type      varchar(10) not null,
    file           varchar(1000),
    product_id     char(32)
        constraint indoormap_api_produc_product_id_8a524e4e_fk_indoormap
            references digimap_db.indoormap_api_product,
    source_url     varchar(255),
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_productattachment
    owner to postgres;

create index idx_21593_indoormap_api_produc_product_id_8a524e4e_fk_indoormap
    on digimap_db.indoormap_api_productattachment (product_id);


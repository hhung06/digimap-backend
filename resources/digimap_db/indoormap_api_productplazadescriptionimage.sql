create table digimap_db.indoormap_api_productplazadescriptionimage
(
    deleted_at     timestamp with time zone,
    restored_at    timestamp with time zone,
    transaction_id char(32),
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    id             char(32)      not null
        constraint idx_21620_primary
            primary key,
    image          varchar(1000) not null,
    "order"        bigint        not null,
    plaza_id       char(32)      not null
        constraint indoormap_api_produc_plaza_id_bd5ae4f5_fk_indoormap
            references digimap_db.indoormap_api_productplaza
);

alter table digimap_db.indoormap_api_productplazadescriptionimage
    owner to postgres;

create index idx_21620_indoormap_a_plaza_i_a9bc82_idx
    on digimap_db.indoormap_api_productplazadescriptionimage (plaza_id);


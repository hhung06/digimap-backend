create table digimap_db.indoormap_api_article
(
    deleted_at             timestamp with time zone,
    restored_at            timestamp with time zone,
    transaction_id         char(32),
    created_at             timestamp with time zone,
    updated_at             timestamp with time zone,
    id                     char(32)     not null
        constraint idx_21246_primary
            primary key,
    external_id            varchar(255),
    title                  varchar(255) not null,
    content                text,
    status                 varchar(50)  not null,
    published_at           timestamp with time zone,
    localization           json,
    product_id             char(32)
        constraint indoormap_api_articl_product_id_e8f3d65f_fk_indoormap
            references digimap_db.indoormap_api_product,
    product_category_id    char(32)
        constraint indoormap_api_articl_product_category_id_395d7522_fk_indoormap
            references digimap_db.indoormap_api_productcategory,
    venue_id               char(32)
        constraint indoormap_api_articl_venue_id_337016f5_fk_indoormap
            references digimap_db.indoormap_api_venue,
    placement              varchar(50)  not null,
    navigate               varchar(50),
    published_period_end   date,
    published_period_start date,
    created_by             varchar(50)  not null,
    location_id            char(32)
        constraint indoormap_api_articl_location_id_aa68488f_fk_indoormap
            references digimap_db.indoormap_api_location,
    application_language   integer,
    company_name           varchar(255),
    email                  varchar(255),
    full_name              varchar(255),
    label                  varchar(255),
    landline               varchar(50),
    mobile                 varchar(50)
);

alter table digimap_db.indoormap_api_article
    owner to postgres;

create index idx_21246_indoormap_api_articl_product_category_id_395d7522_fk_
    on digimap_db.indoormap_api_article (product_category_id);

create index idx_21246_indoormap_api_articl_venue_id_337016f5_fk_indoormap
    on digimap_db.indoormap_api_article (venue_id);

create index idx_21246_indoormap_api_articl_product_id_e8f3d65f_fk_indoormap
    on digimap_db.indoormap_api_article (product_id);

create index idx_21246_indoormap_api_articl_location_id_aa68488f_fk_indoorma
    on digimap_db.indoormap_api_article (location_id);


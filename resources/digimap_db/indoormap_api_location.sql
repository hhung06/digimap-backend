create table digimap_db.indoormap_api_location
(
    created_at                       timestamp with time zone,
    updated_at                       timestamp with time zone,
    deleted_at                       timestamp with time zone,
    id                               char(32)      not null
        constraint idx_21437_primary
            primary key,
    common_name                      varchar(1000) not null,
    externalid                       varchar(255),
    common_description               text,
    common_contact_phone             varchar(255),
    common_social_twitter            varchar(1000),
    common_social_facebook           varchar(1000),
    common_social_website            varchar(1000),
    level_id                         char(32)
        constraint indoormap_api_locati_level_id_efc07011_fk_indoormap
            references digimap_db.indoormap_api_level,
    common_logo                      varchar(1000),
    common_color                     varchar(255),
    place_work_hours                 json,
    venue_id                         char(32)
        constraint indoormap_api_locati_venue_id_ea11932a_fk_indoormap
            references digimap_db.indoormap_api_venue,
    common_location_state            smallint,
    booth_event_date                 date,
    booth_number                     varchar(50),
    booth_products_showcased         text,
    booth_services_offered           text,
    booth_size                       varchar(50),
    common_address                   varchar(255),
    common_contact_email             varchar(255),
    common_hidden                    boolean       not null,
    common_latitude                  numeric(20, 17),
    common_location_state_end_date   date,
    common_location_state_start_date date,
    common_location_type             smallint      not null,
    common_longitude                 numeric(20, 17),
    common_short_name                varchar(1000),
    common_social_instagram          varchar(1000),
    common_social_tiktok             varchar(1000),
    person_full_name                 varchar(255),
    person_job_title                 varchar(255),
    room_bed_count                   integer,
    room_department                  varchar(255),
    room_equipment_details           text,
    room_number                      varchar(50),
    is_top_location                  boolean       not null,
    top_location_sort_index          bigint,
    icon_default                     varchar(255),
    main_category_id                 char(32)
        constraint indoormap_api_locati_main_category_id_62f50b77_fk_indoormap
            references digimap_db.indoormap_api_locationcategory,
    common_large_logo                varchar(1000),
    common_medium_logo               varchar(1000),
    common_small_logo                varchar(1000),
    custom                           json,
    localization                     json,
    source                           varchar(8)    not null,
    top_logo                         varchar(1000),
    common_show_short_name           boolean       not null,
    restored_at                      timestamp with time zone,
    transaction_id                   char(32),
    end_time                         timestamp with time zone,
    start_time                       timestamp with time zone,
    top_logo_type                    varchar(100),
    common_sub_type                  smallint      not null,
    is_searchable                    boolean       not null
);

alter table digimap_db.indoormap_api_location
    owner to postgres;

create index idx_21437_indoormap_api_locati_venue_id_ea11932a_fk_indoormap
    on digimap_db.indoormap_api_location (venue_id);

create index idx_21437_indoormap_api_locati_level_id_efc07011_fk_indoormap
    on digimap_db.indoormap_api_location (level_id);

create index idx_21437_indoormap_api_locati_main_category_id_62f50b77_fk_ind
    on digimap_db.indoormap_api_location (main_category_id);

create index idx_21437_indoormap_a_venue_i_8c464d_idx
    on digimap_db.indoormap_api_location (venue_id);


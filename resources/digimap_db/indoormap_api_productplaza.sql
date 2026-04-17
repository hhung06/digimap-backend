create table digimap_db.indoormap_api_productplaza
(
    deleted_at           timestamp with time zone,
    restored_at          timestamp with time zone,
    transaction_id       char(32),
    created_at           timestamp with time zone,
    updated_at           timestamp with time zone,
    id                   char(32)     not null
        constraint idx_21608_primary
            primary key,
    title                varchar(255) not null,
    description          text,
    status               varchar(50)  not null,
    category             varchar(100),
    country              varchar(100),
    display_type         varchar(50)  not null,
    video                varchar(1000),
    thumbnail            varchar(1000),
    approved_at          timestamp with time zone,
    rejected_at          timestamp with time zone,
    location_id          char(32)     not null
        constraint indoormap_api_produc_location_id_7b7949e7_fk_indoormap
            references digimap_db.indoormap_api_location,
    reviewed_by_id       integer
        constraint indoormap_api_produc_reviewed_by_id_cf864593_fk_auth_user
            references digimap_db.auth_user,
    venue_id             char(32)     not null
        constraint indoormap_api_produc_venue_id_75cd3797_fk_indoormap
            references digimap_db.indoormap_api_venue,
    vimeo_video_id       varchar(100),
    view_count           bigint       not null,
    application_language integer,
    company_name         varchar(255),
    email                varchar(255),
    full_name            varchar(255),
    landline             varchar(50),
    localization         json,
    mobile               varchar(50)
);

alter table digimap_db.indoormap_api_productplaza
    owner to postgres;

create index idx_21608_indoormap_a_locatio_feb59b_idx
    on digimap_db.indoormap_api_productplaza (location_id);

create index idx_21608_indoormap_api_produc_reviewed_by_id_cf864593_fk_auth_
    on digimap_db.indoormap_api_productplaza (reviewed_by_id);

create index idx_21608_indoormap_a_created_5f1d4e_idx
    on digimap_db.indoormap_api_productplaza (created_at);

create index idx_21608_indoormap_a_status_143e2e_idx
    on digimap_db.indoormap_api_productplaza (status);

create index idx_21608_indoormap_a_venue_i_83fbd5_idx
    on digimap_db.indoormap_api_productplaza (venue_id);


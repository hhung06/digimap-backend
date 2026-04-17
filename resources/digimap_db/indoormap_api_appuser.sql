create table digimap_db.indoormap_api_appuser
(
    id               bigserial
        constraint idx_21225_primary
            primary key,
    created_at       timestamp with time zone,
    updated_at       timestamp with time zone,
    deleted_at       timestamp with time zone,
    external_id      varchar(100) not null,
    type             integer,
    section          integer,
    company_name     varchar(255),
    company_name_en  varchar(255),
    department       varchar(255),
    position         varchar(255),
    position_en      varchar(255),
    last_name        varchar(100),
    first_name       varchar(100),
    last_name_en     varchar(100),
    first_name_en    varchar(100),
    tel              varchar(50),
    mail             varchar(254),
    token            text,
    token_expiration timestamp with time zone,
    survey_flag      integer      not null,
    staff_lead_flag  integer      not null,
    app              varchar(100),
    restored_at      timestamp with time zone,
    transaction_id   char(32)
);

alter table digimap_db.indoormap_api_appuser
    owner to postgres;

create index idx_21225_indoormap_a_departm_c8665b_idx
    on digimap_db.indoormap_api_appuser (department);

create index idx_21225_indoormap_a_type_801dc5_idx
    on digimap_db.indoormap_api_appuser (type);

create unique index idx_21225_external_id
    on digimap_db.indoormap_api_appuser (external_id);

create index idx_21225_indoormap_a_section_eefb80_idx
    on digimap_db.indoormap_api_appuser (section);

create index idx_21225_indoormap_a_externa_3d8243_idx
    on digimap_db.indoormap_api_appuser (external_id);

create index idx_21225_indoormap_a_app_7fe2db_idx
    on digimap_db.indoormap_api_appuser (app);


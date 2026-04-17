create table digimap_db.indoormap_api_survey
(
    created_at      timestamp with time zone,
    updated_at      timestamp with time zone,
    deleted_at      timestamp with time zone,
    id              varchar(255) not null
        constraint idx_21712_primary
            primary key,
    created_by_id   bigint
        constraint indoormap_api_survey_created_by_id_e3b2ab61_fk_indoormap
            references digimap_db.indoormap_api_appuser,
    end_date        timestamp with time zone,
    is_forced       boolean      not null,
    publish_type    integer      not null,
    start_date      timestamp with time zone,
    status          integer      not null,
    title           varchar(255),
    app             varchar(20)  not null,
    restored_at     timestamp with time zone,
    transaction_id  char(32),
    content         varchar(500),
    source          integer      not null,
    venue_id        char(32)
        constraint indoormap_api_survey_venue_id_ceb6ac5c_fk_indoormap_api_venue_i
            references digimap_db.indoormap_api_venue,
    segment_filters json,
    external_id     varchar(255)
);

alter table digimap_db.indoormap_api_survey
    owner to postgres;

create index idx_21712_indoormap_api_survey_created_by_id_e3b2ab61_fk_indoor
    on digimap_db.indoormap_api_survey (created_by_id);

create index idx_21712_indoormap_api_survey_venue_id_ceb6ac5c_fk_indoormap_a
    on digimap_db.indoormap_api_survey (venue_id);


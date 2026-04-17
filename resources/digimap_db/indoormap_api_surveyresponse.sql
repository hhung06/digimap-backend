create table digimap_db.indoormap_api_surveyresponse
(
    id             char(32)                 not null
        constraint idx_21731_primary
            primary key,
    submitted_at   timestamp with time zone not null,
    survey_id      varchar(255)             not null
        constraint indoormap_api_survey_survey_id_cf938a9c_fk_indoormap
            references digimap_db.indoormap_api_survey,
    user_id        bigint
        constraint indoormap_api_survey_user_id_c6490791_fk_indoormap
            references digimap_db.indoormap_api_appuser,
    created_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    restored_at    timestamp with time zone,
    transaction_id char(32),
    updated_at     timestamp with time zone
);

alter table digimap_db.indoormap_api_surveyresponse
    owner to postgres;

create index idx_21731_indoormap_api_survey_user_id_c6490791_fk_indoormap
    on digimap_db.indoormap_api_surveyresponse (user_id);

create index idx_21731_indoormap_api_surveyresponse_survey_id_cf938a9c
    on digimap_db.indoormap_api_surveyresponse (survey_id);


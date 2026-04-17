create table digimap_db.indoormap_api_surveyanswer
(
    id             char(32) not null
        constraint idx_21723_primary
            primary key,
    answer_text    text,
    option_id      char(32)
        constraint indoormap_api_survey_option_id_1f032bc1_fk_indoormap
            references digimap_db.indoormap_api_option,
    question_id    char(32) not null
        constraint indoormap_api_survey_question_id_344d0ff0_fk_indoormap
            references digimap_db.indoormap_api_question,
    response_id    char(32) not null
        constraint indoormap_api_survey_response_id_c45e3d07_fk_indoormap
            references digimap_db.indoormap_api_surveyresponse,
    created_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    restored_at    timestamp with time zone,
    transaction_id char(32),
    updated_at     timestamp with time zone
);

alter table digimap_db.indoormap_api_surveyanswer
    owner to postgres;

create index idx_21723_indoormap_api_survey_option_id_1f032bc1_fk_indoormap
    on digimap_db.indoormap_api_surveyanswer (option_id);

create index idx_21723_indoormap_api_survey_question_id_344d0ff0_fk_indoorma
    on digimap_db.indoormap_api_surveyanswer (question_id);

create index idx_21723_indoormap_api_survey_response_id_c45e3d07_fk_indoorma
    on digimap_db.indoormap_api_surveyanswer (response_id);


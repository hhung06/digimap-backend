create table digimap_db.indoormap_api_question
(
    created_at      timestamp with time zone,
    updated_at      timestamp with time zone,
    deleted_at      timestamp with time zone,
    id              char(32)    not null
        constraint idx_21667_primary
            primary key,
    question_number bigint      not null,
    question_type   varchar(20) not null,
    question_text   text        not null,
    survey_id       varchar(255)
        constraint indoormap_api_questi_survey_id_d0aa8913_fk_indoormap
            references digimap_db.indoormap_api_survey,
    is_required     boolean     not null,
    is_other        boolean     not null,
    restored_at     timestamp with time zone,
    transaction_id  char(32)
);

alter table digimap_db.indoormap_api_question
    owner to postgres;

create index idx_21667_indoormap_api_questi_survey_id_d0aa8913_fk_indoormap
    on digimap_db.indoormap_api_question (survey_id);


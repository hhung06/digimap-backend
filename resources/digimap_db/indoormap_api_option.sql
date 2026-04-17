create table digimap_db.indoormap_api_option
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32)     not null
        constraint idx_21552_primary
            primary key,
    option_number  bigint       not null,
    option_text    varchar(255) not null,
    question_id    char(32)
        constraint indoormap_api_option_question_id_32136110_fk_indoormap
            references digimap_db.indoormap_api_question,
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_option
    owner to postgres;

create index idx_21552_indoormap_api_option_question_id_32136110_fk_indoorma
    on digimap_db.indoormap_api_option (question_id);


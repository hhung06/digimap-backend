create table digimap_db.indoormap_api_event_tags
(
    id          bigserial
        constraint idx_21352_primary
            primary key,
    event_id    char(32) not null
        constraint indoormap_api_event__event_id_4582f9f1_fk_indoormap
            references digimap_db.indoormap_api_event,
    eventtag_id char(32) not null
        constraint indoormap_api_event__eventtag_id_60216ac7_fk_indoormap
            references digimap_db.indoormap_api_eventtag
);

alter table digimap_db.indoormap_api_event_tags
    owner to postgres;

create unique index idx_21352_indoormap_api_event_tags_event_id_eventtag_id_9a06894
    on digimap_db.indoormap_api_event_tags (event_id, eventtag_id);

create index idx_21352_indoormap_api_event__eventtag_id_60216ac7_fk_indoorma
    on digimap_db.indoormap_api_event_tags (eventtag_id);


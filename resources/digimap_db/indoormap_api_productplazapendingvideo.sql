create table digimap_db.indoormap_api_productplazapendingvideo
(
    deleted_at     timestamp with time zone,
    restored_at    timestamp with time zone,
    transaction_id char(32),
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    id             char(32)      not null
        constraint idx_21638_primary
            primary key,
    video_id       varchar(100)  not null,
    upload_link    varchar(2000) not null,
    status         varchar(20)   not null
);

alter table digimap_db.indoormap_api_productplazapendingvideo
    owner to postgres;

create unique index idx_21638_video_id
    on digimap_db.indoormap_api_productplazapendingvideo (video_id);


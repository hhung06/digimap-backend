create table digimap_db.indoormap_api_searchquery
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32)     not null
        constraint idx_21688_primary
            primary key,
    search_term    varchar(255) not null,
    search_count   bigint       not null,
    last_searched  timestamp with time zone,
    origin         varchar(255),
    app_id         varchar(255),
    restored_at    timestamp with time zone,
    transaction_id char(32),
    venue_id       varchar(255),
    is_promoted    boolean      not null,
    reference      json,
    status         varchar(20)  not null
);

alter table digimap_db.indoormap_api_searchquery
    owner to postgres;

create index idx_21688_indoormap_api_searchquery_search_term_3379e1da
    on digimap_db.indoormap_api_searchquery (search_term);


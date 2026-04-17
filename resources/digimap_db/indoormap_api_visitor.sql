create table digimap_db.indoormap_api_visitor
(
    deleted_at      timestamp with time zone,
    restored_at     timestamp with time zone,
    transaction_id  char(32),
    created_at      timestamp with time zone,
    updated_at      timestamp with time zone,
    id              char(32) not null
        constraint idx_21824_primary
            primary key,
    full_name       varchar(100),
    phone_number    bytea,
    visitor_type    integer  not null,
    interests       json     not null,
    venue_id        char(32)
        constraint indoormap_api_visito_venue_id_b57e30fd_fk_indoormap
            references digimap_db.indoormap_api_venue,
    email           varchar(255),
    ip_address      varchar(45),
    user_agent      text,
    business_name   varchar(255),
    is_consented    boolean  not null,
    other_interests text
);

alter table digimap_db.indoormap_api_visitor
    owner to postgres;

create index idx_21824_indoormap_api_visito_venue_id_b57e30fd_fk_indoormap
    on digimap_db.indoormap_api_visitor (venue_id);


create table digimap_db.indoormap_api_language
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32)   not null
        constraint idx_21401_primary
            primary key,
    code           varchar(4) not null,
    name           varchar(255),
    date_format    varchar(255),
    venue_id       char(32)
        constraint indoormap_api_langua_venue_id_a5e9c523_fk_indoormap
            references digimap_db.indoormap_api_venue,
    enabled        boolean    not null,
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_language
    owner to postgres;

create index idx_21401_indoormap_a_enabled_a51d3f_idx
    on digimap_db.indoormap_api_language (enabled);

create index idx_21401_indoormap_api_language_venue_id_a5e9c523
    on digimap_db.indoormap_api_language (venue_id);

create index idx_21401_indoormap_a_code_199921_idx
    on digimap_db.indoormap_api_language (code);


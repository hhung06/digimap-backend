create table digimap_db.indoormap_api_venuetheme
(
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    deleted_at     timestamp with time zone,
    id             char(32) not null
        constraint idx_21803_primary
            primary key,
    venue_id       char(32),
    data           json,
    restored_at    timestamp with time zone,
    transaction_id char(32)
);

alter table digimap_db.indoormap_api_venuetheme
    owner to postgres;

create unique index idx_21803_indoormap_api_venuetheme_venue_id_108e010f_uniq
    on digimap_db.indoormap_api_venuetheme (venue_id);


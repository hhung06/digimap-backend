create table digimap_db.indoormap_api_venue
(
    id             char(32)         not null
        constraint idx_21773_primary
            primary key,
    name           varchar(200)     not null,
    created_at     timestamp with time zone,
    updated_at     timestamp with time zone,
    address        varchar(255),
    city           varchar(255),
    countrycode    varchar(255),
    defaultmap     varchar(255),
    deleted_at     timestamp with time zone,
    externalid     varchar(255)     not null,
    latitude       double precision not null,
    longitude      double precision not null,
    postal         varchar(255),
    telephone      varchar(255),
    timezone       varchar(255),
    largelogo      varchar(1000),
    mediumlogo     varchar(1000),
    originallogo   varchar(1000),
    smalllogo      varchar(1000),
    publish        boolean          not null,
    customer_id    char(32)
        constraint indoormap_api_venue_customer_id_0de13a77_fk_indoormap
            references digimap_db.indoormap_api_customer,
    country        varchar(150),
    state          varchar(150),
    workhours      text             not null,
    type           smallint         not null,
    description    text,
    public_key     varchar(43)      not null,
    custom_data    json,
    localization   json,
    private_key    varchar(64)      not null,
    endat          timestamp with time zone,
    restored_at    timestamp with time zone,
    seodescription text,
    seokeywords    varchar(255),
    seotitle       varchar(255),
    startat        timestamp with time zone,
    subdomains     text,
    transaction_id char(32),
    appconfigs     json,
    appdomains     json,
    bodytag        text,
    headtag        text
);

alter table digimap_db.indoormap_api_venue
    owner to postgres;

create unique index idx_21773_unique_externalid
    on digimap_db.indoormap_api_venue (externalid);

create index idx_21773_indoormap_api_venue_customer_id_0de13a77_fk_indoormap
    on digimap_db.indoormap_api_venue (customer_id);

create unique index idx_21773_indoormap_api_venue_externalid_c85249c5_uniq
    on digimap_db.indoormap_api_venue (externalid);


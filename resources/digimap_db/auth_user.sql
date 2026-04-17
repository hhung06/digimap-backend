create table digimap_db.auth_user
(
    id           serial
        constraint idx_21110_primary
            primary key,
    password     varchar(128)             not null,
    last_login   timestamp with time zone,
    is_superuser boolean                  not null,
    username     varchar(150)             not null,
    first_name   varchar(150)             not null,
    last_name    varchar(150)             not null,
    email        varchar(254)             not null,
    is_staff     boolean                  not null,
    is_active    boolean                  not null,
    date_joined  timestamp with time zone not null
);

alter table digimap_db.auth_user
    owner to postgres;

create unique index idx_21110_username
    on digimap_db.auth_user (username);


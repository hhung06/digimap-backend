create table digimap_db.indoormap_api_perspective
(
    created_at            timestamp with time zone,
    updated_at            timestamp with time zone,
    deleted_at            timestamp with time zone,
    id                    char(32)         not null
        constraint idx_21558_primary
            primary key,
    name                  varchar(255),
    camerazoom            double precision not null,
    cameratype            smallint         not null,
    cameramaxzoom         double precision not null,
    cameraminzoom         double precision not null,
    cameratargetbearing   double precision not null,
    cameratargetcenterlat double precision not null,
    cameratargetcenterlng double precision not null,
    cameratargetpitch     double precision not null,
    cameratargetzoom      double precision not null,
    restored_at           timestamp with time zone,
    transaction_id        char(32)
);

alter table digimap_db.indoormap_api_perspective
    owner to postgres;


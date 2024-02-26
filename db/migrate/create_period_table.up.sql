create table period
(
    id          bigserial
        primary key,
    name        varchar(255)          not null,
    description text                  not null,
    year        integer               not null,
    begin_at    date                  not null,
    ended_at    date                  not null,
    created_at  timestamp(0),
    updated_at  timestamp(0),
    deleted     boolean default false not null
);
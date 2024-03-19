DROP TABLE IF EXISTS period;

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

DROP TABLE IF EXISTS users;
create table users
(
    id                 bigserial
        primary key,
    name               varchar(255)                               not null,
    email              varchar(255)                               not null
        constraint users_email_unique
            unique,
    mac_address        macaddr                                    not null
        constraint users_mac_address_unique
            unique,
    telegram_id        varchar(255)
        constraint users_telegram_id_unique
            unique,
    email_verified_at  timestamp(0),
    password           varchar(255)                               not null,
    remember_token     varchar(100),
    device_mac_address varchar(17)
        constraint users_device_mac_address_unique
            unique,
    created_at         timestamp(0),
    updated_at         timestamp(0),
    deleted            boolean      default false                 not null,
    department_id      integer      default 0                     not null,
    role               varchar(255) default ''::character varying not null,
    surname            varchar(255) default ''::character varying not null,
    lastname           varchar(255) default ''::character varying not null,
    birth_date         date                                       not null,
    position           varchar(255) default ''::character varying
);

insert into public.users (id, name, email, mac_address, telegram_id, email_verified_at, password, remember_token, device_mac_address, created_at, updated_at, deleted, department_id, role, surname, lastname, birth_date, position)
values  (57, 'Иван', 'test1@gmail.com', 'd0:2b:20:eb:3d:9d', '1466023675', null, '$2a$14$9jltXblz65ai/j9PTxJBC.s8ubTsYciXp7GiAQl0aOP.NZXhbeVY.', null, null, '2021-04-04 19:07:24', null, false, 1, '', 'Иван', 'Иван', '2021-04-16', ''),
        (4, 'Тимофей', 'test2@mail.ru', 'bc:e1:43:49:fe:cc', '284632731', null, 'drsdf', null, null, null, null, false, 0, 'Employee', '', '', '2021-04-16', ''),
        (1, 'Анастасия', 'test3@officetime.tech', '7c:a1:ae:a9:c0:d8', '343536263', null, '$2y$10$dxogytXGfNkv.9Bg4veOsu/Bksb/11hbNtQPSh6068pxq4UoVVgdC', '2EBv5lByUt7wsIibjCxKNmm7niDOvIUykoijBuKoSSNpNeqSrFc1vYcknh6s', null, null, null, false, 1, 'Admin', 'Некрасова', 'Николаевна', '1990-09-30', 'Backend Developer');

DROP TABLE IF EXISTS router;

create table router
(
    id          serial
        constraint router_pk
            primary key,
    uuid
    name        varchar(255) default ''::character varying not null,
    description varchar(255) default ''::character varying not null,
    address     varchar(20)  default ''::character varying not null,
    login       varchar(50)  default ''::character varying not null,
    password    varchar(255) default ''::character varying not null,
    created_at  timestamp                                  not null,
    updated_at  timestamp,
    status      boolean      default false,
    work_time   boolean      default false                 not null
);

create unique index router_address_uindex
    on router (address);

create unique index router_name_uindex
    on router (name);

insert into public.router (id, name, description, address, login, password, created_at, updated_at, status, work_time)
values  (1, 'Router1', 'Коридор', '95.84.134.115:8728', 'admin', 'Vtlcgjgek1', '2021-04-07 19:29:35.000000', null, true, true);
create table customers
(
    id          uuid    not null,
    name        varchar not null,
    email       varchar,
    created_at  date,
    modified_at date,
    constraint customers_pk
        primary key (id),
    constraint customers_pk_2
        unique (email)
);
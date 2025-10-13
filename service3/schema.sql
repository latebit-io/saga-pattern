create table loans
(
    id                  uuid      not null,
    customer_id         uuid      not null,
    mortgage_id         uuid      not null,
    loan_amount         numeric   not null,
    interest_rate       numeric   not null,
    term_years          int       not null,
    monthly_payment     numeric   not null,
    outstanding_balance numeric   not null,
    status              varchar   not null,
    start_date          timestamp not null,
    maturity_date       timestamp not null,
    created_at          timestamp not null,
    modified_at         timestamp not null,
    constraint loans_pk
        primary key (id)
);

create table payments
(
    id               uuid      not null,
    loan_id          uuid      not null,
    customer_id      uuid      not null,
    payment_amount   numeric   not null,
    principal_amount numeric   not null,
    interest_amount  numeric   not null,
    payment_date     timestamp not null,
    payment_type     varchar   not null,
    created_at       timestamp not null,
    constraint payments_pk
        primary key (id)
);
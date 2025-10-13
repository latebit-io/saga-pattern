create table mortgage_applications
(
    id              uuid      not null,
    customer_id     uuid      not null,
    loan_amount     numeric   not null,
    property_value  numeric   not null,
    interest_rate   numeric   not null,
    term_years      int       not null,
    status          varchar   not null,
    created_at      timestamp not null,
    modified_at     timestamp not null,
    constraint mortgage_applications_pk
        primary key (id)
);
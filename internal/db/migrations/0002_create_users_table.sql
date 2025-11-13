create table if not exists users
(
    id         serial primary key unique,
    name       varchar(255)        not null,
    email      varchar(255) unique not null,
    password   varchar(255)        not null,
    created_at timestamp
)
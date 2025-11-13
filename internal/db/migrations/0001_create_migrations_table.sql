create table if not exists migrations
(
    id         serial primary key,
    name       varchar(255) unique,
    created_at timestamp
)
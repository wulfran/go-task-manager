create table if not exists tasks
(
    id serial primary key unique,
    name varchar(64) not null,
    priority int not null,
    description varchar(255),
    due_date timestamp,
    created_at timestamp,
    created_by int,
    constraint fk_created_by foreign key (created_by) references users(id)
)
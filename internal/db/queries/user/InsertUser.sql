-- InsertUser
insert into users (name, email, password, created_at)
values ($1, $2, $3, $4)
returning id
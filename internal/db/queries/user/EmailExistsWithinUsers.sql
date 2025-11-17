-- EmailExistsWithinUsers
select exists(
    select id
    from users
    where email=$1
)

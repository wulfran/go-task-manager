select exists(
    select id
    from tasks
    where created_by = $1
    and id = $2
)
update tasks
set name = $1,
    priority = $2,
    description = $3,
    due_date = $4
where id = $5
select name, priority, description, due_date,created_at, created_by
from tasks
where id = $1
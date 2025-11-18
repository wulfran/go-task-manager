select id, name, priority, description, due_date, created_at, created_by
from tasks
where created_by=$1
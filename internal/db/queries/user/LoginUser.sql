select id, name, email, password, created_at
from users
where email=$1

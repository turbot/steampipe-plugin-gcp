select self_link, instance_name, kind, type
from gcp.gcp_sql_backup
where id = {{ output.backup_id.value }};
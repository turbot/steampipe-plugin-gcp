select id, instance_name
from gcp.gcp_sql_backup
where id = {{ output.backup_id.value }}0 and instance_name = '{{ resourceName }}';
select name, database_version, backup_enabled
from gcp.gcp_sql_database_instance
where title = '{{ resourceName }}';
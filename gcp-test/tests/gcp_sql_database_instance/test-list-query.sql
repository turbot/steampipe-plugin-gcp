select name, database_version
from gcp.gcp_sql_database_instance
where title = '{{ resourceName }}';
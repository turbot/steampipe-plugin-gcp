select tags, title, akas
from gcp.gcp_sql_database_instance
where name = '{{ resourceName }}';
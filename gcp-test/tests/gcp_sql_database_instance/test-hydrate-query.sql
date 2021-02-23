select name, jsonb_array_elements(instance_users) -> 'host' as host_name
from gcp.gcp_sql_database_instance
where name = '{{ resourceName }}';
select name, database_version, machine_type
from gcp.gcp_sql_database_instance
where name = 'dummy-{{ resourceName }}';
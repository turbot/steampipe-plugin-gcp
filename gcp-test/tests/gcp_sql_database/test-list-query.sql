select name, instance_name
from gcp.gcp_sql_database
where name = '{{ resourceName }}'
select name, instance_name, kind
from gcp.gcp_sql_database
where name = 'dummy-{{ resourceName }}' and instance_name = 'dummy-{{ resourceName }}'
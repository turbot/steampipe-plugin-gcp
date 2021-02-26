select name, instance_name, kind
from gcp.gcp_sql_database
where name = '{{ resourceName }}' and instance_name = '{{ resourceName }}';
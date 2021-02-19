select name, instance_name, kind, charset, self_link, location, project
from gcp.gcp_sql_database
where name = '{{ resourceName }}' and instance_name = '{{ resourceName }}';
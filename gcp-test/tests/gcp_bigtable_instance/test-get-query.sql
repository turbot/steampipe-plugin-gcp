select name, display_name, instance_type, location, project
from gcp.gcp_bigtable_instance
where name = '{{ resourceName }}';
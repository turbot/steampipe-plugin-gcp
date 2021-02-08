select name, display_name, instance_type, state
from gcp.gcp_bigtable_instance
where name = 'dummy-{{ resourceName }}'
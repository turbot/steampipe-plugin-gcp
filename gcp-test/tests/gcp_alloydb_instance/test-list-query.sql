select cluster_name, instance_display_name
from gcp.gcp_alloydb_instance
where cluster_name = '{{ resourceName }}'
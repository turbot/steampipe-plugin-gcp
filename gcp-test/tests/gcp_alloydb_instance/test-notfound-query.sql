select cluster_name, title, akas
from gcp.gcp_alloydb_instance
where cluster_name = '{{ resourceName }}:asdf'
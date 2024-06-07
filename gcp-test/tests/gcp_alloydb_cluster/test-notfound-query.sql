select name, title, akas
from gcp.gcp_alloydb_cluster
where display_name = '{{ resourceName }}:asdf'
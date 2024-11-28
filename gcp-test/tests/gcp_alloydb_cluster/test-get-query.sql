select name, display_name
from gcp.gcp_alloydb_cluster
where display_name = '{{ resourceName }}';
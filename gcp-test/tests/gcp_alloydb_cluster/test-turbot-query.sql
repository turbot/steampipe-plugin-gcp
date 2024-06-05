select title, akas
from gcp.gcp_alloydb_instance
where display_name = '{{ resourceName }}';
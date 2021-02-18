select name, display_name
from gcp.gcp_bigtable_instance
where title = '{{ resourceName }}';
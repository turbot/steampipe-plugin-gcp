select tags, title, akas
from gcp.gcp_bigtable_instance
where name = '{{ resourceName }}';
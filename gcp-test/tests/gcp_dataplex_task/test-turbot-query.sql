select title, akas
from gcp.gcp_dataplex_task
where display_name = '{{ resourceName }}';
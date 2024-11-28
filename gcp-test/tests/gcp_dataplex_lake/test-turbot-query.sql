select title, akas
from gcp.gcp_dataplex_lake
where display_name = '{{ resourceName }}';
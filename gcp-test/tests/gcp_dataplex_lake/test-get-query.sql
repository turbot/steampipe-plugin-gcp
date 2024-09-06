select name, display_name
from gcp.gcp_dataplex_lake
where display_name = '{{ resourceName }}';
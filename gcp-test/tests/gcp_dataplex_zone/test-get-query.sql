select name, display_name, lake_name
from gcp.gcp_dataplex_zone
where name = '{{ output.resource_id.value }}';
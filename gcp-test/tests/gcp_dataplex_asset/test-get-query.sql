select name, display_name, lake_name, zone_name
from gcp.gcp_dataplex_asset
where name = '{{ output.resource_id.value }}';
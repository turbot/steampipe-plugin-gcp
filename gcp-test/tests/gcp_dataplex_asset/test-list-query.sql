select name, display_name
from gcp.gcp_dataplex_asset
where zone_name = '{{ output.zone_name.value }}'
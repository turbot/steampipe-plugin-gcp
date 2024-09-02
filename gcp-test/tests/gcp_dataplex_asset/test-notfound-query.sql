select name, title, akas
from gcp.gcp_dataplex_asset
where zone_name = '{{ output.zone_name.value }}abd'
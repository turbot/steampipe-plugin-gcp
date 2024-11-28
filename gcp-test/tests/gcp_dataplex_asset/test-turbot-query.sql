select title, akas
from gcp.gcp_dataplex_asset
where zone_name = '{{ output.zone_name.value }}' and display_name = '{{ resourceName }}';
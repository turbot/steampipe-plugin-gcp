select name, display_name
from gcp.gcp_dataplex_zone
where akas::text = '["{{ output.resource_aka.value }}"]'
select name, display_name
from gcp.gcp_dataplex_task
where akas::text = '["{{ output.resource_aka.value }}"]'
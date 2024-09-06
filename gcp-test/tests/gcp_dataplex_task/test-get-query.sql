select name, display_name, lake_name
from gcp.gcp_dataplex_task
where name = '{{ output.resource_id.value }}';
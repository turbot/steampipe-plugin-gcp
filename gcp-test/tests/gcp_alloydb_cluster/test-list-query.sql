select name, display_name
from gcp.gcp_alloydb_cluster
where akas::text = '["{{ output.resource_aka.value }}"]'
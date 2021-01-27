select name, title, akas, tags
from gcp.gcp_storage_bucket
where akas::text = '["{{ output.resource_aka.value }}"]'
select name, title, akas, labels, tags
from gcp.gcp_storage_bucket
where akas::text = '["{{ output.resource_aka.value }}"]'
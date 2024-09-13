select name, project, title, akas
from gcp.gcp_dataproc_metastore_service
where akas::text = '["{{ output.resource_aka.value }}"]'
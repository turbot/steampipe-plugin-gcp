select title, akas
from gcp.gcp_dataproc_metastore_service
where name = '{{ output.resource_id.value }}';
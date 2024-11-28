select name, title, state, uid, project
from gcp.gcp_dataproc_metastore_service
where name = '{{ output.resource_id.value }}';
select cluster_name, project, title, akas
from gcp.gcp_dataproc_cluster
where akas::text = '["{{ output.resource_aka.value }}"]'
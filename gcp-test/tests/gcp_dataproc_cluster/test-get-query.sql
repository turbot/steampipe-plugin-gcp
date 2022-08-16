select cluster_name, title, project
from gcp.gcp_dataproc_cluster
where cluster_name = '{{ resourceName }}';
select cluster_name, title, akas
from gcp.gcp_dataproc_cluster
where cluster_name = '{{ resourceName }}:asdf'
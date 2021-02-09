select name, iam_policy 
from gcp.gcp_compute_subnetwork 
where name = '{{ resourceName }}'
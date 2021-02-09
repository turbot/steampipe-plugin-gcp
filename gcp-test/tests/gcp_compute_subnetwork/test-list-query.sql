select 
  name, 
  description, 
  kind
from 
  gcp.gcp_compute_subnetwork
where 
  title = '{{ resourceName }}'
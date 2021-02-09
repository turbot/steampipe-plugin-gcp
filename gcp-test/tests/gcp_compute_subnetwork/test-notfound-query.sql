select 
  name,
  id,
  description, 
  kind
from 
  gcp.gcp_compute_subnetwork
where 
  name = 'dummy-{{ resourceName }}'
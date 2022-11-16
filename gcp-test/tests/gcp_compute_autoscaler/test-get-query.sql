select 
  name,
  description,
  kind,
  self_link, 
  status,
  location, 
  project,
  recommended_size
from 
  gcp_compute_autoscaler
where 
  name = '{{resourceName}}';
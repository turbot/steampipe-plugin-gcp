select 
  name,
  description,
  self_link, 
  size,
  title,
  akas, 
  location, 
  project
from 
  gcp_compute_instance_group
where 
  name = '{{resourceName}}-dummy'
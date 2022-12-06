select 
  name,
  description,
  self_link, 
  size::text,
  title,
  akas, 
  location, 
  project
from 
  gcp_compute_instance_group
where 
  akas::text = '["{{ output.resource_aka.value }}"]';
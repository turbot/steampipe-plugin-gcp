select 
  title,
  akas
from 
  gcp_compute_autoscaler
where 
  name = '{{resourceName}}';

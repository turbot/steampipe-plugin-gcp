select
  name,
  akas,
  title
from
  gcp.gcp_compute_subnetwork
where
  name = '{{ resourceName }}'
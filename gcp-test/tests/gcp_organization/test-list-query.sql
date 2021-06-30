select
  name,
  organization_id,
  lifecycle_state,
  display_name,
  title,
  akas
from
  gcp.gcp_organization
where
  name = '{{ output.name.value }}';
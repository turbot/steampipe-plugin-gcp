select
  name,
  title
from
  gcp_project
where
  name = 'dummy-{{ resourceName }}';
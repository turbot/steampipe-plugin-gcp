select project, location, title, akas
from gcp.gcp_project_organization_policy
where title = '{{ output.resource_title.value }}';
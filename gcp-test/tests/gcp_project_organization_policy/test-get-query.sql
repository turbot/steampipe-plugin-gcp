select title, akas, project
from gcp.gcp_project_organization_policy
where id = '{{ output.resource_id.value }}';
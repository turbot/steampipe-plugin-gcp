select id, project, title, akas
from gcp.gcp_projects_organization_policy
where title = '{{ output.resource_title.value }}:asdf';
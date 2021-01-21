select name, role_id, title, stage, deleted, description, included_permissions
from gcp.gcp_iam_role
where name = '{{ output.resource_id.value }}'
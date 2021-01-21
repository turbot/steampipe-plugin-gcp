select title, akas
from gcp.gcp_iam_role
where name = '{{ output.resource_id.value }}'
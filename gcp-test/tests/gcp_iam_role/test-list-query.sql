select name, role_id
from gcp.gcp_iam_role
where title = '{{resourceName}}'
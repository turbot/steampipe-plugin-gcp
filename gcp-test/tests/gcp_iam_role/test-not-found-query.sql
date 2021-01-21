select name, deleted, stage, role_id
from gcp.gcp_iam_role
where name = 'dummy-{{resourceName}}'
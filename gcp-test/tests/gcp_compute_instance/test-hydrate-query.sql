select name, iam_policy
from gcp.gcp_compute_instance
where name = '{{ resourceName }}';
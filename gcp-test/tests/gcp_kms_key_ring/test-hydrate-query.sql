select name, iam_policy
from gcp.gcp_kms_key_ring
where name = '{{ resourceName }}'
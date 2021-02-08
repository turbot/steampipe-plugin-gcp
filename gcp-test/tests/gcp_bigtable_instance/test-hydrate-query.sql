select name, iam_policy
from gcp.gcp_bigtable_instance
where name = '{{ resourceName }}'
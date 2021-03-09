select name, iam_policy
from gcp.gcp_compute_disk
where name = '{{ resourceName }}:asdfghjkl'
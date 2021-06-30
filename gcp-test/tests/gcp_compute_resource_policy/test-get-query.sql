select name, description, status, self_link, kind, instance_schedule_policy
from gcp.gcp_compute_resource_policy
where name = '{{ resourceName }}';
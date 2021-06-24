select name, self_link, description
from gcp.gcp_logging_bucket
where name = 'dummy-{{resourceName}}' and location = '{{ output.region_id.value }}';
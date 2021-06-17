select name, retention_days, location, project
from gcp.gcp_logging_bucket
where name = '{{ output.resource_id.value }}';
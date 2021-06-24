select name, retention_days, location, project
from gcp.gcp_logging_bucket
where name = '{{ resourceName }}' and location = '{{ output.region_id.value }}';
select title, akas, location, project
from gcp.gcp_logging_bucket
where name = '{{ output.resource_id.value }}';
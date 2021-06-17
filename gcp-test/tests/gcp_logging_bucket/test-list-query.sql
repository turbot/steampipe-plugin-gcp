select name, title
from gcp.gcp_logging_bucket
where name = '{{ output.resource_id.value }}'
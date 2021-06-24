select name, title
from gcp.gcp_logging_bucket
where name = '{{ resourceName }}';
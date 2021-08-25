select name, destination
from gcp.gcp_logging_sink
where name = '{{resourceName}}';
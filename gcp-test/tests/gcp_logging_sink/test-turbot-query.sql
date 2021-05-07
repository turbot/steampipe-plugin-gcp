select title, akas
from gcp.gcp_logging_sink
where name = '{{resourceName}}';
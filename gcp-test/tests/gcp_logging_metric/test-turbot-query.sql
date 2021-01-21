select title, akas
from gcp.gcp_logging_metric
where name = '{{resourceName}}'
select name, description
from gcp.gcp_logging_metric
where name = '{{resourceName}}'
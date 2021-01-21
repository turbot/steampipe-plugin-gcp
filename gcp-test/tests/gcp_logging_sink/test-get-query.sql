select name, destination, disabled, filter, description, include_children
from gcp.gcp_logging_sink
where name = '{{resourceName}}'
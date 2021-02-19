select name, destination, disabled, filter, description, include_children, project, location
from gcp.gcp_logging_sink
where name = '{{ resourceName }}';
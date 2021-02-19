select name, disabled, filter, description, location, project
from gcp.gcp_logging_exclusion
where name = '{{ resourceName }}';
select name, disabled, description, filter
from gcp.gcp_logging_exclusion
where name = 'projects/parker-aaa/locations/global/dummy-{{resourceName}}'
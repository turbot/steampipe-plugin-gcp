select name, disabled, filter, description
from gcp.gcp_logging_exclusion
where name = '{{resourceName}}'
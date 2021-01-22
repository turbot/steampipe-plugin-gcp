select name, display_name
from gcp.gcp_monitoring_group
where display_name = '{{resourceName}}'
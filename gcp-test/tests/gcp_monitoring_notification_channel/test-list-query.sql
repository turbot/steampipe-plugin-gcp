select name, display_name
from gcp.gcp_monitoring_notification_channel
where display_name = '{{resourceName}}'
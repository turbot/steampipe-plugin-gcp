select name, display_name, enabled, description
from gcp.gcp_monitoring_notification_channel
where name = 'dummy-{{resourceName}}'
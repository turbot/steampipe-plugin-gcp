select name, display_name, enabled, type, description, labels, project, location
from gcp.gcp_monitoring_notification_channel
where name = '{{ output.resource_id.value.split("/").pop() }}';
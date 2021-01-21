select title, tags, akas
from gcp.gcp_monitoring_notification_channel
where name = '{{ output.resource_id.value.split("/").pop() }}'
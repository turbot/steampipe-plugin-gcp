select name, display_name, combiner, enabled, tags, project
from gcp.gcp_monitoring_alert_policy
where name = '{{ output.resource_id.value.split("/").pop() }}'
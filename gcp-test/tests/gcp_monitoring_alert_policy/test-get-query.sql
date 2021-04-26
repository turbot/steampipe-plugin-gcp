select name, display_name, combiner, enabled, documentation, project, location, user_labels
from gcp.gcp_monitoring_alert_policy
where name = '{{ output.resource_id.value.split("/").pop() }}';
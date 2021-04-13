select name, display_name, combiner
from gcp.gcp_monitoring_alert_policy
where display_name = '{{ resourceName }}';
select name, display_name, filter, is_cluster, parent_name, project
from gcp.gcp_monitoring_group
where name = '{{ output.resource_id.value.split("/").pop() }}'
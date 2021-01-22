select title, akas
from gcp.gcp_monitoring_group
where name = '{{ output.resource_id.value.split("/").pop() }}'
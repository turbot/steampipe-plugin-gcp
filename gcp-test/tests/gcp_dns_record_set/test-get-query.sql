select name, managedzone_name, type, ttl, title, location, project
from gcp.gcp_dns_record_set
where name = '{{ output.record_set_name.value }}' and type = '{{ output.resource_type.value }}';
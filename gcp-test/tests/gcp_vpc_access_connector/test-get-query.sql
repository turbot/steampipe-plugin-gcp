select name, location, state, network
from gcp.gcp_vpc_access_connector
where name = '{{ output.resource_id.value }}';
select name, state
from gcp.gcp_vpc_access_connector
where name = '{{ output.resource_id.value }}-dummy';
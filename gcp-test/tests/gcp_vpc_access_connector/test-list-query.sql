select name, title
from gcp.gcp_vpc_access_connector
where akas::text = '["{{ output.resource_aka.value }}"]';
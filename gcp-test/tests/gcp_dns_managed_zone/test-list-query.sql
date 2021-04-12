select name, description
from gcp.gcp_dns_managed_zone
where akas::text = '["{{ output.resource_aka.value }}"]';
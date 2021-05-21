select name, managedzone_name, type
from gcp.gcp_dns_record_set
where akas::text = '["{{ output.resource_aka.value }}"]';
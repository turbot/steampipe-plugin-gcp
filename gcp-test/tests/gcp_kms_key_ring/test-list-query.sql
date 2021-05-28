select name
from gcp.gcp_kms_key_ring
where akas::text = '["{{ output.resource_aka.value }}"]'
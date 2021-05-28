select name, rotation_period
from gcp.gcp_kms_key
where akas::text = '["{{ output.resource_aka.value }}"]';
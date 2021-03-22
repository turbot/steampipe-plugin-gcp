select name
from gcp.gcp_kms_key
where name = 'dummy-{{output.resource_id.value}}'
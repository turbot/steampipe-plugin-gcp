select name
from gcp.gcp_kms_key
where name = 'dummy-{{ resourceName }}';
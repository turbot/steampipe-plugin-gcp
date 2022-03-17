select name
from gcp.gcp_kms_key_version
where name = 'dummy-{{ resourceName }}';
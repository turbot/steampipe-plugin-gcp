select key_name
from gcp.gcp_kms_key_version
where key_name = 'dummy-{{ resourceName }}';
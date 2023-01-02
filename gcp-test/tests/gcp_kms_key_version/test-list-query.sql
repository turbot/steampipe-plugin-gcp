select key_name, crypto_key_version, self_link, akas 
from gcp.gcp_kms_key_version 
where key_name = '{{ resourceName }}';
select name, crypto_key_versions, self_link, akas 
from gcp.gcp_kms_key_version 
where name = '{{ resourceName }}';
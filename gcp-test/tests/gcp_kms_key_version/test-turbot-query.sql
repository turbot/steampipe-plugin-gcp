select name, crypto_key_versions, self_link
from gcp.gcp_kms_key_version
where name = '{{ resourceName }}';
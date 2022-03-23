select name, crypto_key_version, self_link 
from gcp.gcp_kms_key_version 
where name = '{{ resourceName }}' and key_ring_name = '{{ resourceName }}' and location = 'global' and crypto_key_version = 1;

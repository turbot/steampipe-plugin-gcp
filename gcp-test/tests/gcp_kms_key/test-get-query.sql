select name, rotation_period
from gcp.gcp_kms_key
where name = '{{ resourceName }}' and key_ring_name = '{{ resourceName }}' and location = 'global';
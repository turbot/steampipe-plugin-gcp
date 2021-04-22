select name
from gcp.gcp_kms_key_ring
where name = 'dummy-{{ resourceName }}'
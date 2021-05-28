select title, akas
from gcp.gcp_kms_key_ring
where name = '{{ resourceName }}'
select name, id, address, address_type
from gcp_compute_global_address
where name = 'dummy-{{resourceName}}'
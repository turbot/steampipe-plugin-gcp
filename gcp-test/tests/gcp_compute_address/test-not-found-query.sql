select name, id, address, address_type
from gcp.gcp_compute_address
where name = 'dummy-{{resourceName}}'
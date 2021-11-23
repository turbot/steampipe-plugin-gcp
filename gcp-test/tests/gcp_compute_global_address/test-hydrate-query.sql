select name, address_type, description, kind, network_tier, self_link
from gcp_compute_global_address
where name = '{{resourceName}}'
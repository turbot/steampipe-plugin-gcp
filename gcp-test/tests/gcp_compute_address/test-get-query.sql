select name, address_type, description, kind, network_tier, self_link, location, project, subnetwork
from gcp.gcp_compute_address
where name = '{{resourceName}}';
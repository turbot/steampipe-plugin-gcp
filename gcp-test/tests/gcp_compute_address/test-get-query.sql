select name, address_type, description, kind, network_tier, self_link, region, project, subnetwork
from gcp.gcp_compute_address
where name = '{{resourceName}}'
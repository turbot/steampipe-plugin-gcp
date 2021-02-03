select name, id, kind, description, region, self_link, forwarding_rules, project
from gcp.google_compute_vpn_gateway
where name = '{{ resourceName }}'
select name, description, network_name, kind, description, enable_flow_logs, log_config_enable, ip_cidr_range, network, region, self_link, secondary_ip_ranges, project, location 
from gcp.gcp_compute_subnetwork 
where name = '{{ resourceName }}'
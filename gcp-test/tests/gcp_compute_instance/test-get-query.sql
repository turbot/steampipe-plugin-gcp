select name, description, machine_type_name, can_ip_forward, cpu_platform, deletion_protection, kind, label_fingerprint, labels, machine_type, service_accounts, network_tags, zone, zone_name, location, project, self_link
from gcp.gcp_compute_instance
where name = '{{ resourceName }}';
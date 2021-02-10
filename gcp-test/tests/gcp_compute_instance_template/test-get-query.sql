select name, description, instance_can_ip_forward, instance_description, instance_disks, instance_machine_type, instance_metadata, instance_network_interfaces, instance_scheduling, instance_service_accounts, instance_tags, kind, location, project, self_link
from gcp.gcp_compute_instance_template
where name = '{{ resourceName }}'
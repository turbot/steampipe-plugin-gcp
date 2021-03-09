select name, description, kind, self_link, project, maintenance_policy, autoscaling_policy_mode, autoscaling_policy_max_nodes, autoscaling_policy_min_nodes
from gcp.gcp_compute_node_group
where name = '{{ resourceName }}'
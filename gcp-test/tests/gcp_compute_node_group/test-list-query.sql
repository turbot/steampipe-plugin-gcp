select name, autoscaling_policy_max_nodes, size
from gcp.gcp_compute_node_group
where title = '{{ resourceName }}'
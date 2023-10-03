select title, akas
from gcp.gcp_compute_machine_type
where name = '{{ output.machine_type.value }}' and zone = 'us-east1-b';
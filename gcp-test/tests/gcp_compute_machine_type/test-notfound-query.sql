select name, id
from gcp.gcp_compute_machine_type
where name = '{{ output.machine_type.value }}.dd'
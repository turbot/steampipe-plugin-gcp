select name, title, akas, kind, guest_cpus, memory_mb, image_space_gb, maximum_persistent_disks, maximum_persistent_disks_size_gb, is_shared_cpu
from gcp.gcp_compute_machine_type
where name = '{{ output.machine_type.value }}';
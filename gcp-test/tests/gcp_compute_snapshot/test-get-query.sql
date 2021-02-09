select name, title, auto_created, description, disk_size_gb, source_disk, akas, project
from gcp.gcp_compute_snapshot
where name = '{{ resourceName }}'
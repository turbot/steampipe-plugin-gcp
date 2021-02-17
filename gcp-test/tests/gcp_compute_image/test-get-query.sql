select name, kind, description, self_link, disk_size_gb, source_disk, project, labels
from gcp.gcp_compute_image
where name = '{{ resourceName }}' and project = '{{ output.project_id.value }}';
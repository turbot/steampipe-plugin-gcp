select name, kind, description, disk_size_gb, self_link, source_disk, iam_policy
from gcp.gcp_compute_image
where name = '{{ resourceName }}' and project = '{{ output.project_id.value }}';
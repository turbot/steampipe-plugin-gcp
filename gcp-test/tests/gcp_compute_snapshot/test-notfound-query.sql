select name, description, source_disk, self_link, project, title, akas
from gcp.gcp_compute_snapshot
where name = '{{ resourceName }}:asdf'
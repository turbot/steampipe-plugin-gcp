select name, description, kind, enable_cdn, self_link, bucket_name, project
from gcp.gcp_compute_backend_bucket
where name = '{{ resourceName }}'
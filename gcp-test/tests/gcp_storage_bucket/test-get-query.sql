select name, id, self_link, title, akas, tags
from gcp.gcp_storage_bucket
where name = '{{ resourceName }}'
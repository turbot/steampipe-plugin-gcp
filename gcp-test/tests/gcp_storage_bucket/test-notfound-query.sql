select name, self_link, title, akas
from gcp.gcp_storage_bucket
where name = '["{{ output.resource_aka.value }}"]:asd'
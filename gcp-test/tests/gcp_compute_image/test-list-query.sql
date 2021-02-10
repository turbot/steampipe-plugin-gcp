select name, description
from gcp.gcp_compute_image
where title = '{{ resourceName }}'
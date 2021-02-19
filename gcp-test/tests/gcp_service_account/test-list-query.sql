select name, display_name
from gcp.gcp_service_account
where display_name = '{{ resourceName }}'
select name, unique_id, display_name, email, disabled, description, oauth2_client_id, location, project
from gcp.gcp_service_account
where name = '{{ output.name.value }}';
select name, default_service, description, kind, self_link, location, project, location_type, tests, path_matchers, host_rules
from gcp.gcp_compute_url_map
where name = '{{ resourceName }}'
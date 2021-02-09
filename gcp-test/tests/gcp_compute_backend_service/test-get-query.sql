select name, enable_cdn, description, kind, load_balancing_scheme, self_link, affinity_cookie_ttl_sec, project, location, location_type, connection_draining_timeout_sec, log_config_enable, port, port_name, protocol, signed_url_cache_max_age_sec, health_checks
from gcp.gcp_compute_backend_service
where name = '{{ resourceName }}'
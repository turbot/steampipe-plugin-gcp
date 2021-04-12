select name, description, kind, authorization_policy, proxy_bind, quic_override, self_link, url_map, ssl_certificates, location_type, project, location
from gcp.gcp_compute_target_https_proxy
where name = '{{ resourceName }}';
select name, description, kind, proxy_header, service, self_link, ssl_certificates, project, location
from gcp.gcp_compute_target_ssl_proxy
where name = '{{ resourceName }}';
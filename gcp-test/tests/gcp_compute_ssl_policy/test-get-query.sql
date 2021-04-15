select name, description, kind, fingerprint, min_tls_version, profile, enabled_features, self_link, project, location
from gcp.gcp_compute_ssl_policy
where name = '{{ resourceName }}';
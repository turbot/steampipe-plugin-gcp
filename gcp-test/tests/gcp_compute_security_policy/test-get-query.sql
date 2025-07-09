select name, description, self_link, fingerprint, type, rules, project, location
from gcp_compute_security_policy
where name = '{{ resourceName }}';
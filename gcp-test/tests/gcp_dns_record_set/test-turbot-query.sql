select title, akas
from gcp.gcp_dns_record_set
where name = '{{ output.record_set_name.value }}' and type = 'A';
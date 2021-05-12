select title, tags, akas
from gcp.gcp_bigquery_table
where id = '{{ output.project_id.value }}:{{ resourceName }}.{{ resourceName }}';
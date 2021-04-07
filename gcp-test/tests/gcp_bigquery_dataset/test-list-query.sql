select name, description
from gcp.gcp_bigquery_dataset
where id = '{{ output.project_id.value }}:{{ resourceName }}';
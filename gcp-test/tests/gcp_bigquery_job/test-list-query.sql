select job_id, id, kind
from gcp.gcp_bigquery_job
where id = '{{ output.resource_id.value }}';
select job_id, id, kind, project, location
from gcp.gcp_bigquery_job
where job_id = '{{ resourceName }}';
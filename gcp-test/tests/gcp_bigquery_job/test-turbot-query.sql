select title, akas
from gcp.gcp_bigquery_job
where job_id = '{{ resourceName }}';
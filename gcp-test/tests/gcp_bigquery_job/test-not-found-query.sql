select id, kind
from gcp.gcp_bigquery_job
where job_id = 'dummy-{{ resourceName }}';
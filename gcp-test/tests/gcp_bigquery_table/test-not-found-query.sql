select dataset_id, id, kind
from gcp.gcp_bigquery_table
where id = 'dummy-{{ resourceName }}';
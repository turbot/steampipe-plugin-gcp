select name, id, kind, description
from gcp.gcp_bigquery_dataset
where dataset_id = 'dummy-{{ resourceName }}';
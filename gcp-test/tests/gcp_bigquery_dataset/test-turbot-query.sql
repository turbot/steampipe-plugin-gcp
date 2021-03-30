select title, tags, akas
from gcp.gcp_bigquery_dataset
where dataset_id = '{{ resourceName }}';
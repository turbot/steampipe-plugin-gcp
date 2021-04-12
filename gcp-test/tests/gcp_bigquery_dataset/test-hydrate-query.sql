select dataset_id, description, kind, etag, self_link
from gcp.gcp_bigquery_dataset
where dataset_id = '{{ resourceName }}';
select name, dataset_id, id, description, kind, etag, default_partition_expiration_ms, default_table_expiration_ms, self_link, access, labels, project, location
from gcp.gcp_bigquery_dataset
where dataset_id = '{{ resourceName }}';
select dataset_id, id
from gcp.gcp_bigquery_table
where dataset_id = '{{ output.parent_resource_id.value }}';
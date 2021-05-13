select 
  dataset_id, 
  table_id, 
  kind
from 
  gcp.gcp_bigquery_table
where 
  table_id = 'dummy_{{ resourceName }}' 
  and dataset_id = 'dummy_{{ resourceName }}';

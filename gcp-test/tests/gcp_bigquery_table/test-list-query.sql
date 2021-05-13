select 
  dataset_id, 
  table_id
from 
  gcp.gcp_bigquery_table
where 
  table_id = '{{ resourceName }}';

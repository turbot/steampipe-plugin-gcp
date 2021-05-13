select 
  dataset_id, 
  id
from 
  gcp.gcp_bigquery_table
where 
  table_id = '' 
  and dataset_id = '{{ resourceName }}';

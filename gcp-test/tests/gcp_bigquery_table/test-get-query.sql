select 
  table_id,
  dataset_id, 
  id, 
  kind, 
  labels, 
  range_partitioning, 
  time_partitioning, 
  type, 
  self_link,  
  title, 
  location, 
  project
from 
  gcp.gcp_bigquery_table
where 
  table_id = '{{ resourceName }}'
  and dataset_id = '{{ resourceName }}';

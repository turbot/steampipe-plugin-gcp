select 
  title, 
  tags, 
  akas
from 
  gcp.gcp_bigquery_table
where 
  table_id = '{{ resourceName }}'
  and dataset_id = '{{ resourceName }}';

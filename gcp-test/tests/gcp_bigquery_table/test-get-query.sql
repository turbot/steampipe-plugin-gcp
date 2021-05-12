select dataset_id, friendly_name, id, kind, labels, range_partitioning, table_reference, time_partitioning, type, self_link, view, title, location, project
from gcp.gcp_bigquery_table
where id='{{ output.project_id.value }}:{{ output.parent_resource_id.value }}.{{ output.parent_resource_id.value }}';

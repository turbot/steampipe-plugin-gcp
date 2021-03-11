select title, akas
from gcp.gcp_sql_database
where name = '{{ resourceName }}';
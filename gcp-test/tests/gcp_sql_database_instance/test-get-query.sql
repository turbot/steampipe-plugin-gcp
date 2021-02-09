select name, instance_type, database_version, kind, machine_type, data_disk_type, data_disk_size_gb, self_link, storage_auto_resize, enable_point_in_time_recovery, crash_safe_replication_enabled, connection_name, location, project
from gcp.gcp_sql_database_instance
where name = '{{ resourceName }}'
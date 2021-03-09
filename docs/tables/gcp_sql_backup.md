# Table: gcp_sql_backup

Backups help to restore lost data to Cloud SQL instance. Backups protect data from loss or damage.

## Examples

### Basic info

```sql
select
  id,
  instance_name,
  description,
  status,
  end_time,
  enqueued_time,
  start_time,
  window_start_time
from
  gcp_sql_backup;
```


### Count of backups by their type (i.e AUTOMATED and ON_DEMAND)

```sql
select
  type,
  count(*) as backup_count
from
  gcp_sql_backup
group by
  type;
```


### Get the error message if the backup failed

```sql
select
  id,
  instance_name,
  e ->> 'code' as error_code,
  e ->> 'message' as error_message
from
  gcp_sql_backup,
  jsonb_array_elements(error) as e
where
  status = 'FAILED';
```
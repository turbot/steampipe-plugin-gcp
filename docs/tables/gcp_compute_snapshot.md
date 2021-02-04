# Table: gcp_compute_snapshot

Snapshot is the backup of disk. It is a global resource, so it can be use to restore data to a new disk or instance within the same project.

## Examples

### Count of snapshots per disk

```sql
select
  source_disk_name,
  count(*) as snapshot_count
from
  gcp_compute_snapshot
group by
  source_disk_name;
```

### List of manually created snapshots

```sql
select
  name,
  source_disk_name,
  auto_created
from
  gcp_compute_snapshot
where
  not auto_created;
```

### Disk info for each snapshot

```sql
select
  s.name as snapshot_name,
  d.name as disk_name,
  d.size_gb as disk_size,
  d.type_name as disk_type,
  d.location_type
from
  gcp_compute_snapshot as s
  join gcp_compute_disk as d on s.source_disk = d.self_link;
```

### List of snapshots older than 90 days

```sql
select
  name,
  creation_timestamp,
  age(creation_timestamp)
from
  gcp_compute_snapshot
where
  creation_timestamp <= (current_date - interval '90' day)
order by
  creation_timestamp;
```

### List of snapshots with Google-managed key

```sql
select
  name,
  source_disk,
  self_link
from
  gcp_compute_snapshot
where
  kms_key_name is null;
```

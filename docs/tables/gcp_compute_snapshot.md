# Table: gcp_compute_snapshot

Snapshot is the backup of disk. It is a global resource, so it can be use to restore data to a new disk or instance within the same project.

## Examples

### Count of snapshots per disk

```sql
select
split_part(source_disk, '/', 11) as source_disk,
count(*) as snapshot_count
from
  gcp_compute_snapshot
group by
  (source_disk);
```


### List of manual created snapshots

```sql
select
  name,
  split_part(source_disk, '/', 11) as source_disk,
  auto_created
from
  gcp_compute_snapshot
where
  not auto_created;
```


### Disk info of each snapshot

```sql
select
  s.name as snapshot_name,
  d.name as disk_name,
  d.size_gb as disk_size,
  split_part(d.type, '/', 11) as disk_type,
  d.location_type
from
  gcp_compute_snapshot as s
  join gcp_compute_disk as d on s.source_disk = d.self_link;
```
# Table: gcp_compute_instance

Compute Engine manages the physical disks and the data distribution to ensure redundancy and optimal performance. Persistent disks are located independently from virtual machine (VM) instances, so it can detach or move persistent disks to keep data even after deletion of instances.

## Examples

### Basic info

```sql
select
  name,
  id,
  size_gb as disk_size_in_gb,
  type_name,
  zone_name,
  region_name,
  location_type
from
  gcp_compute_disk;
```

### List of disks with Google-managed key

```sql
select
  name,
  id,
  zone_name,
  disk_encryption_kms_key
from
  gcp_compute_disk
where
  disk_encryption_kms_key is null;
```

### List of disks that are not in use

```sql
select
  name,
  id,
  users
from
  gcp_compute_disk
where
  users is null;
```

### List of disks that are Regional

```sql
select
  name,
  region_name
from
  gcp_compute_disk
where
  location_type = 'REGIONAL';
```

### Disk count in each availability zone

```sql
select
  zone_name,
  count(*)
from
  gcp_compute_disk
group by
  zone_name
order by
  count desc;
```

### List Disks by size

```sql
select
  name,
  size_gb
from
  gcp_compute_disk
order by
  size_gb desc;
```

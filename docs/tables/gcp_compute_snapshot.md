---
title: "Steampipe Table: gcp_compute_snapshot - Query Compute Engine Snapshots using SQL"
description: "Allows users to query Compute Engine Snapshots in Google Cloud Platform (GCP), specifically the snapshot's ID, name, description, status, and other metadata, providing insights into disk snapshot usage and management."
folder: "Compute"
---

# Table: gcp_compute_snapshot - Query Compute Engine Snapshots using SQL

A Compute Engine Snapshot in Google Cloud Platform (GCP) is a copy of a virtual machine's disk at a specific point in time. Snapshots are used for backing up data, duplicating disks, and creating images for boot disks. They capture the entire state of a disk, including all data, installed applications, and system settings.

## Table Usage Guide

The `gcp_compute_snapshot` table provides insights into Compute Engine Snapshots within Google Cloud Platform (GCP). As a system administrator, explore snapshot-specific details through this table, including snapshot status, storage locations, and associated metadata. Utilize it to uncover information about snapshots, such as their creation time, disk size, and the source disk's ID.

## Examples

### Count of snapshots per disk
Analyze the distribution of snapshots across various disks to understand the backup frequency for each disk. This can be useful for ensuring regular backups are being made for all important data.

```sql+postgres
select
  source_disk_name,
  count(*) as snapshot_count
from
  gcp_compute_snapshot
group by
  source_disk_name;
```

```sql+sqlite
select
  source_disk_name,
  count(*) as snapshot_count
from
  gcp_compute_snapshot
group by
  source_disk_name;
```

### List of manually created snapshots
Explore which disk snapshots in your Google Cloud Platform (GCP) Compute Engine have been manually created. This can help you manage and differentiate them from automatically generated snapshots.

```sql+postgres
select
  name,
  source_disk_name,
  auto_created
from
  gcp_compute_snapshot
where
  not auto_created;
```

```sql+sqlite
select
  name,
  source_disk_name,
  auto_created
from
  gcp_compute_snapshot
where
  auto_created = 0;
```

### Disk info for each snapshot
Explore the details of each snapshot in your Google Cloud Platform's compute engine, including the associated disk's name, size, type, and location. This is particularly useful for managing resources and understanding the composition of your snapshots.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which system snapshots have been stored for over 90 days. This can be helpful for managing storage resources and identifying potential data that could be archived or deleted.

```sql+postgres
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

```sql+sqlite
select
  name,
  creation_timestamp,
  julianday('now') - julianday(creation_timestamp) as age
from
  gcp_compute_snapshot
where
  julianday(creation_timestamp) <= julianday('now', '-90 day')
order by
  creation_timestamp;
```

### List of snapshots with Google-managed key
Explore which snapshots are using Google-managed keys. This is useful to ensure your data is properly encrypted and secure.

```sql+postgres
select
  name,
  source_disk,
  self_link
from
  gcp_compute_snapshot
where
  kms_key_name is null;
```

```sql+sqlite
select
  name,
  source_disk,
  self_link
from
  gcp_compute_snapshot
where
  kms_key_name is null;
```
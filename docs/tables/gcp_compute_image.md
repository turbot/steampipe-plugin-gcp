# Table: gcp_compute_image

An Image resource contains a boot loader, an operating system and a root file system that is necessary for starting an instance.

### Compute image basic info

```sql
select
  name,
  id,
  kind,
  status
from
  gcp_compute_image;
```

### List of standard compute images

```sql
select
  name,
  id,
  source_project
from
  gcp_compute_image
where
  deprecated is null;
```

### List of compute images which are not encrypted

```sql
select
  name,
  id,
  image_encryption_key
from
  gcp_compute_image
where
  image_encryption_key is null;
```

### List of compute images which do not have owner tag key

```sql
select
  name,
  id
from
  gcp_compute_image
where
  tags -> 'owner' is null;
```

### List of compute images older than 90 days

```sql
select
  name,
  creation_timestamp,
  age(creation_timestamp)
from
  gcp_compute_image
where
  creation_timestamp <= (current_date - interval '90' day)
order by
  creation_timestamp;
```

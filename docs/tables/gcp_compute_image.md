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


### Get the source disk information

```sql
select
  name,
  source_disk,
  source_disk_id
from
  gcp_compute_image;
```


### List images with policies

```sql
select
  name,
  iam_policy
from
  gcp_compute_image;
```
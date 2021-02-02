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
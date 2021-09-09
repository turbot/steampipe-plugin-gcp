# Table: gcp_compute_image

An Image resource contains a boot loader, an operating system and a root file system that is necessary for starting an instance.

## Examples

### Compute image basic info

```sql
select
  name,
  id,
  kind,
  status,
  deprecation_state
from
  gcp_compute_image;
```

### List of active, standard compute images

```sql
select
  name,
  id,
  source_project,
  deprecation_state
from
  gcp_compute_image
where
  deprecation_state is null
  and is_standard_image;
```

### List of custom (user-defined) images defined in this project

```sql
select
  name,
  id,
  source_project
from
  gcp_compute_image
where
  not is_standard_image;
```

### List of custom (user-defined) images which are not encrypted with a customer key

```sql
select
  name,
  id,
  image_encryption_key
from
  gcp_compute_image
where
  image_encryption_key is null and
  not is_standard_image;
```

### List of user-defined compute images which do not have owner tag key

```sql
select
  name,
  id,
  tags
from
  gcp_morales_aaa.gcp_compute_image
where  tags -> 'owner' is null
  and source_project = project
```

### List of active compute images older than 90 days

```sql
select
  name,
  creation_timestamp,
  age(creation_timestamp),
  deprecation_state
from
  gcp_compute_image
where
  creation_timestamp <= (current_date - interval '90' day)
  and deprecation_state = 'ACTIVE'
order by
  creation_timestamp;
```

### Find VM instances built from images older than 90 days

```sql
select
  vm.name as instance_name,
  d.name as disk_name,
  img.name as image,
  img.creation_timestamp as image_creation_time,
  age(img.creation_timestamp) as image_age,
  img.deprecation_state
from
  gcp_compute_instance as vm,
  jsonb_array_elements(vm.disks) as vmd,
  gcp_compute_disk as d,
  gcp_compute_image as img
where
  vmd ->> 'source' = d.self_link
  and (vmd ->> 'boot') :: bool
  and d.source_image = img.self_link
  and img.creation_timestamp <= (current_date - interval '90' day);
```

### Find VM instances built from deprecated, deleted, or obsolete images

```sql
select
  vm.name as instance_name,
  d.name as disk_name,
  img.name as image,
  img.creation_timestamp as image_creation_time,
  age(img.creation_timestamp) as image_age,
  img.deprecation_state
from
  gcp_compute_instance as vm,
  jsonb_array_elements(vm.disks) as vmd,
  gcp_compute_disk as d,
  gcp_compute_image as img
where
  vmd ->> 'source' = d.self_link
  and (vmd ->> 'boot') :: bool
  and d.source_image = img.self_link
  and deprecation_state != 'ACTIVE';
```

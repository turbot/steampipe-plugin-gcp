---
title: "Steampipe Table: gcp_compute_image - Query Google Cloud Compute Engine Images using SQL"
description: "Allows users to query Google Cloud Compute Engine Images, providing detailed information about each image, including its creation timestamp, description, disk size, and more."
folder: "Compute"
---

# Table: gcp_compute_image - Query Google Cloud Compute Engine Images using SQL

Google Cloud Compute Engine Images are binary data that is used to create instances in Google Cloud Compute Engine. These images contain a boot loader, an operating system, and a root file system. Google Cloud Compute Engine Images are essential for creating and managing instances in Google Cloud Compute Engine.

## Table Usage Guide

The `gcp_compute_image` table provides insights into images within Google Cloud Compute Engine. As a Cloud Engineer, you can explore image-specific details through this table, including their creation timestamps, descriptions, disk sizes, and more. Utilize it to uncover information about images, such as those that are deprecated, the operating systems they contain, and their source disk IDs.

## Examples

### Compute image basic info
Explore the basic information of your Google Cloud Platform compute images to identify their status and deprecation state. This can assist in managing and maintaining your compute resources effectively.

```sql+postgres
select
  name,
  id,
  kind,
  status,
  deprecation_state
from
  gcp_compute_image;
```

```sql+sqlite
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
Explore the active compute images that are sourced from different projects. This can help in managing resources and identifying potential redundancies.

```sql+postgres
select
  name,
  id,
  source_project
from
  gcp_compute_image
where
  deprecation_state = 'ACTIVE'
  and source_project != project;
```

```sql+sqlite
select
  name,
  id,
  source_project
from
  gcp_compute_image
where
  deprecation_state = 'ACTIVE'
  and source_project != project;
```

### List of custom (user-defined) images defined in this project
Explore which custom images have been defined within a specific project. This can help in understanding the customization and modifications made to the project.

```sql+postgres
select
  name,
  id,
  source_project
from
  gcp_compute_image
where
  source_project = project;
```

```sql+sqlite
select
  name,
  id,
  source_project
from
  gcp_compute_image
where
  source_project = project;
```

### List of compute images which are not encrypted with a customer key
Explore which compute images are not secured with a unique customer key. This query is useful to identify potential security vulnerabilities in your GCP compute images.

```sql+postgres
select
  name,
  id,
  image_encryption_key
from
  gcp_compute_image
where
  image_encryption_key is null;
```

```sql+sqlite
select
  name,
  id,
  image_encryption_key
from
  gcp_compute_image
where
  image_encryption_key is null;
```

### List of user-defined compute images which do not have owner tag key
Explore which user-defined compute images lack an owner tag key. This is useful to identify potential gaps in image management, ensuring that all images are properly attributed to an owner.

```sql+postgres
select
  name,
  id
from
  gcp_compute_image
where
  tags -> 'owner' is null
  and  source_project = project;
```

```sql+sqlite
select
  name,
  id
from
  gcp_compute_image
where
  json_extract(tags, '$.owner') is null
  and  source_project = project;
```

### List of active compute images older than 90 days
Explore which compute images have remained active for more than 90 days. This can help identify areas for potential optimization or cleanup in your GCP environment.

```sql+postgres
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

```sql+sqlite
select
  name,
  creation_timestamp,
  julianday('now') - julianday(creation_timestamp) as age,
  deprecation_state
from
  gcp_compute_image
where
  julianday(creation_timestamp) <= julianday(datetime('now','-90 day'))
  and deprecation_state = 'ACTIVE'
order by
  creation_timestamp;
```

### Find VM instances built from images older than 90 days
This query is useful for maintaining the security and effectiveness of your virtual machine instances. It helps identify any instances that were built from images older than 90 days, allowing you to update or replace them as necessary to ensure optimal performance and compliance with best practices.

```sql+postgres
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

```sql+sqlite
select
  vm.name as instance_name,
  d.name as disk_name,
  img.name as image,
  img.creation_timestamp as image_creation_time,
  julianday('now') - julianday(img.creation_timestamp) as image_age,
  img.deprecation_state
from
  gcp_compute_instance as vm,
  json_each(vm.disks) as vmd,
  gcp_compute_disk as d,
  gcp_compute_image as img
where
  json_extract(vmd.value, '$.source') = d.self_link
  and json_extract(vmd.value, '$.boot') = 1
  and d.source_image = img.self_link
  and img.creation_timestamp <= date('now','-90 day');
```

### Find VM instances built from deprecated, deleted, or obsolete images
Determine the instances where virtual machines (VMs) are built using outdated, deleted, or obsolete images. This is useful for identifying potential security risks and ensuring optimal performance by keeping your VMs up-to-date.

```sql+postgres
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

```sql+sqlite
select
  vm.name as instance_name,
  d.name as disk_name,
  img.name as image,
  img.creation_timestamp as image_creation_time,
  strftime('%s', 'now') - strftime('%s', img.creation_timestamp) as image_age,
  img.deprecation_state
from
  gcp_compute_instance as vm,
  json_each(vm.disks) as vmd,
  gcp_compute_disk as d,
  gcp_compute_image as img
where
  json_extract(vmd.value, '$.source') = d.self_link
  and (json_extract(vmd.value, '$.boot') = 'true')
  and d.source_image = img.self_link
  and deprecation_state != 'ACTIVE';
```
---
title: "Steampipe Table: gcp_compute_instance_template - Query Google Cloud Compute Engine Instance Templates using SQL"
description: "Allows users to query Instance Templates in Google Cloud Compute Engine, specifically providing details about the configuration of virtual machine instances in an instance group."
folder: "Compute"
---

# Table: gcp_compute_instance_template - Query Google Cloud Compute Engine Instance Templates using SQL

An Instance Template in Google Cloud Compute Engine is a resource that provides a way to create instances based on a template. This allows for the uniform deployment of multiple instances that have the same configurations. Instance Templates define the machine type, boot disk image or container image, labels, and other instance properties.

## Table Usage Guide

The `gcp_compute_instance_template` table provides insights into Instance Templates within Google Cloud Compute Engine. As a DevOps engineer, explore template-specific details through this table, including machine type, boot disk image, labels, and other instance properties. Utilize it to uncover information about templates, such as those with specific configurations, the uniform deployment of multiple instances, and the verification of instance properties.

## Examples

### List of c2-standard-4 machine type instance template
Discover the segments that utilize the 'c2-standard-4' machine type within Google Cloud Platform's compute instances. This is particularly useful for assessing resource allocation and planning for future infrastructure needs.

```sql+postgres
select
  name,
  id,
  instance_machine_type
from
  gcp_compute_instance_template
where
  instance_machine_type = 'c2-standard-4';
```

```sql+sqlite
select
  name,
  id,
  instance_machine_type
from
  gcp_compute_instance_template
where
  instance_machine_type = 'c2-standard-4';
```

### Boot Disk info of each instance template
Determine the characteristics of each instance template's boot disk in a GCP compute environment. This can be useful to assess the disk type, size, and source image, which can aid in capacity planning and performance optimization.
```sql+postgres
select
  name,
  id,
  disk ->> 'deviceName' as disk_device_name,
  disk -> 'initializeParams' ->> 'diskType' as disk_type,
  disk -> 'initializeParams' ->> 'diskSizeGb' as disk_size_gb,
  split_part(
    disk -> 'initializeParams' ->> 'sourceImage',
    '/',
    5
  ) as source_image,
  disk ->> 'mode' as mode
from
  gcp_compute_instance_template,
  jsonb_array_elements(instance_disks) as disk;
```

```sql+sqlite
Error: SQLite does not support split functions.
```

### List of SPECIFIC_RESERVATION Instance type instance template
Analyze the settings to understand which instance templates in the Google Cloud Platform are specifically set to consume reservations. This can help optimize resource allocation and cost management in cloud environments.
```sql+postgres
select
  name,
  id,
  instance_reservation_affinity ->> 'consumeReservationType' as consume_reservation_type
from
  gcp_compute_instance_template
where
  instance_reservation_affinity ->> 'consumeReservationType' = 'SPECIFIC_RESERVATION';
```

```sql+sqlite
select
  name,
  id,
  json_extract(instance_reservation_affinity, '$.consumeReservationType') as consume_reservation_type
from
  gcp_compute_instance_template
where
  json_extract(instance_reservation_affinity, '$.consumeReservationType') = 'SPECIFIC_RESERVATION';
```

### Network interface info of each instance template
Determine the network interface details of each instance template in your GCP Compute Engine to understand the network configuration and access settings in each instance. This will help in identifying any irregularities or inconsistencies in the network setup.

```sql+postgres
select
  name,
  id,
  i ->> 'name' as name,
  split_part(i ->> 'network', '/', 10) as network_name,
  p ->> 'name' as access_config_name,
  p ->> 'networkTier' as access_config_network_tier,
  p ->> 'type' as access_config_type
from
  gcp_compute_instance_template,
  jsonb_array_elements(instance_network_interfaces) as i,
  jsonb_array_elements(i -> 'accessConfigs') as p;
```

```sql+sqlite
Error: SQLite does not support split or string_to_array functions.
```

### List of instance templates where instance_can_ip_forward is true
Discover the segments that have the ability to forward IP, a useful feature for routing network traffic effectively and securely. This can be particularly beneficial in scenarios where you need to manage traffic flow across different network interfaces.

```sql+postgres
select
  name,
  id,
  instance_can_ip_forward
from
  gcp_compute_instance_template
where
  instance_can_ip_forward;
```

```sql+sqlite
select
  name,
  id,
  instance_can_ip_forward
from
  gcp_compute_instance_template
where
  instance_can_ip_forward = 1;
```
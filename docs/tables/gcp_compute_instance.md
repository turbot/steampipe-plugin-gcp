---
title: "Steampipe Table: gcp_compute_instance - Query Google Compute Engine Instances using SQL"
description: "Allows users to query Google Compute Engine Instances, providing detailed information about each instance's configuration, status, and associated resources."
folder: "Compute"
---

# Table: gcp_compute_instance - Query Google Compute Engine Instances using SQL

Google Compute Engine is a service within Google Cloud Platform that provides scalable, high-performance virtual machines (VMs) that run on Google's innovative data center infrastructure. It offers a flexible computing environment that lets you choose from pre-defined or custom machine types to suit your workloads. Google Compute Engine also allows you to use various operating systems, frameworks, and languages to build your applications.

## Table Usage Guide

The `gcp_compute_instance` table provides comprehensive information about instances within Google Compute Engine. As a system administrator or DevOps engineer, you can leverage this table to explore instance-specific details, including machine type, operating system, network configuration, and status. This table is particularly useful for assessing the state of your instances, identifying instances with non-standard configurations, and understanding the distribution of resources across your instances.

## Examples

### Instance count in each availability zone
Analyze the distribution of instances across different availability zones to understand the load balancing within your Google Cloud Platform. This can help optimize resource allocation and improve system resilience.

```sql+postgres
select
  zone_name,
  count(*)
from
  gcp_compute_instance
group by
  zone_name
order by
  count desc;
```

```sql+sqlite
select
  zone_name,
  count(*)
from
  gcp_compute_instance
group by
  zone_name
order by
  count(*) desc;
```

### Count the number of instances by instance type
Identify the distribution of different machine types within your Google Cloud Compute instances. This allows you to understand which machine types are most commonly used, aiding in resource allocation and cost management.

```sql+postgres
select
  machine_type_name,
  count(*) as count
from
  gcp_compute_instance
group by
  machine_type_name
order by
  count desc;
```

```sql+sqlite
select
  machine_type_name,
  count(*) as count
from
  gcp_compute_instance
group by
  machine_type_name
order by
  count desc;
```

### List of instances without application label
Identify instances where the application label is missing. This is useful for maintaining consistent tagging practices across your Google Cloud Compute instances.

```sql+postgres
select
  name,
  tags
from
  gcp_compute_instance
where
  tags -> 'application' is null;
```

```sql+sqlite
select
  name,
  tags
from
  gcp_compute_instance
where
  json_extract(tags, '$.application') is null;
```

### List instances having deletion protection feature disabled
Determine the areas in which instances lack deletion protection, a feature that safeguards against accidental data loss. This query is useful for identifying potential vulnerabilities and ensuring data integrity.

```sql+postgres
select
  name,
  deletion_protection
from
  gcp_compute_instance
where
  not deletion_protection;
```

```sql+sqlite
select
  name,
  deletion_protection
from
  gcp_compute_instance
where
  deletion_protection = 0;
```

### List the disk stats attached to the instances
Determine the overall storage capacity across all instances by analyzing the number of disks attached and their respective sizes. This aids in managing resources and planning for storage expansion.

```sql+postgres
select
  name,
  count(d) as num_disks,
  sum( (d ->> 'diskSizeGb') :: int ) as total_storage
from
  gcp_compute_instance as i,
  jsonb_array_elements(disks) as d
group by
  name;
```

```sql+sqlite
select
  i.name,
  count(d.value) as num_disks,
  sum(json_extract(d.value, '$.diskSizeGb')) as total_storage
from
  gcp_compute_instance i,
  json_each(i.disks) as d
group by
  i.name;
```

### Find instances with IP in a given CIDR range
Identify instances within a specific IP range in your Google Cloud Platform's compute instances. This can be useful in understanding your network distribution and identifying potential security risks.

```sql+postgres
select
  name,
  nic ->> 'networkIP' as ip_address
from
  gcp_compute_instance as i,
  jsonb_array_elements(network_interfaces) as nic
where
  (nic ->> 'networkIP') :: inet <<= '10.128.0.0/16' ;
```

```sql+sqlite
select
  i.name,
  json_extract(nic.value, '$.networkIP') as ip_address
from
  gcp_compute_instance i,
  json_each(i.network_interfaces) as nic
where
  json_extract(nic.value, '$.networkIP') like '10.128.%';
```

### Find instances that have been stopped for more than 30 days
Explore instances that have been inactive for an extended period of time. This is useful for identifying potential cost-saving opportunities by eliminating unused resources.

```sql+postgres
select
  name,
  status,
  last_stop_timestamp
from
  gcp_compute_instance
where
  status = 'TERMINATED'
  and last_stop_timestamp < current_timestamp - interval '30 days' ;
```

```sql+sqlite
select
  name,
  status,
  last_stop_timestamp
from
  gcp_compute_instance
where
  status = 'TERMINATED'
  and last_stop_timestamp < datetime('now', '-30 days');
```

### Find the boot disk of each instance
This query allows you to identify the boot disk associated with each virtual machine instance, particularly useful when you need to understand the source image used for the boot disk, such as in instances where you're troubleshooting or auditing system configurations.

```sql+postgres
select
  vm.name as instance_name,
  d.name as disk_name,
  d.source_image
from
  gcp_compute_instance as vm,
  jsonb_array_elements(vm.disks) as vmd,
  gcp_compute_disk as d
where
  vmd ->> 'source' = d.self_link
  and (vmd ->> 'boot') :: bool
  and d.source_image like '%debian-10-buster-v20201014';
```

```sql+sqlite
select
  vm.name as instance_name,
  d.name as disk_name,
  d.source_image
from
  gcp_compute_instance as vm,
  json_each(vm.disks) as vmd,
  gcp_compute_disk as d
where
  json_extract(vmd.value, '$.source') = d.self_link
  and json_extract(vmd.value, '$.boot') = 'true'
  and d.source_image like '%debian-10-buster-v20201014';
```
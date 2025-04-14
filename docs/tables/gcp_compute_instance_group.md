---
title: "Steampipe Table: gcp_compute_instance_group - Query Google Cloud Compute Engine Instance Groups using SQL"
description: "Allows users to query Google Cloud Compute Engine Instance Groups, providing insights into the configuration, status, and properties of these groups."
folder: "Compute"
---

# Table: gcp_compute_instance_group - Query Google Cloud Compute Engine Instance Groups using SQL

Google Cloud Compute Engine Instance Groups are collections of virtual machine (VM) instances that you can manage as a single entity. When they are managed by Instance Group Managers, these groups are called [Managed Instance Groups (MIG)](https://cloud.google.com/compute/docs/instance-groups#managed_instance_groups), and are ideal for highly available applications that require a lot of computing power and need to scale rapidly to meet demand. They offer a range of features including autoscaling, autohealing, regional (multiple zone) deployment, and automatic updating. Otherwise, these groups are called [Unmanaged Instance Groups](https://cloud.google.com/compute/docs/instance-groups#unmanaged_instance_groups), and can contain heterogeneous instances that you can arbitrarily add and remove from them, but do not offer autoscaling, autohealing, rolling update support, multi-zone support, or the use of instance templates and are not a good fit for deploying highly available and scalable workloads, they can just be used for load balancing.

## Table Usage Guide

The `gcp_compute_instance_group` table provides insights into instance groups within Google Cloud Compute Engine. As a system administrator, you can explore group-specific details through this table, including configuration, associated instances, and autoscaling policies. Utilize it to monitor the status of your instance groups, manage load balancing, and plan for capacity adjustments.

## Examples

### Basic Info
Discover the segments of your Google Cloud Platform (GCP) that contain instance groups, gaining insights into aspects like size and location. This can help in project management and resource allocation within the GCP infrastructure.

```sql+postgres
select
  name,
  description,
  self_link,
  size,
  location,
  akas,
  project
from
  gcp_compute_instance_group;
```

```sql+sqlite
select
  name,
  description,
  self_link,
  size,
  location,
  akas,
  project
from
  gcp_compute_instance_group;
```

### Get number of instances per instance group
Analyze the distribution of instances across different groups in your Google Cloud Platform's compute engine. This allows you to manage resources effectively and plan for future scaling needs.

```sql+postgres
select
  name,
  size as no_of_instances
from
  gcp_compute_instance_group;
```

```sql+sqlite
select
  name,
  size as no_of_instances
from
  gcp_compute_instance_group;
```

### Get instance details of each instance group
Explore the status of each instance within a group to gain insights into their operational status. This can help in managing resources more effectively and identifying any instances that may be experiencing issues.

```sql+postgres
select
  g.name,
  ins.name as instance_name,
  ins.status as instance_status
from
  gcp_compute_instance_group as g,
  jsonb_array_elements(instances) as i,
  gcp_compute_instance as ins
where
  (i ->> 'instance') = ins.self_link;
```

```sql+sqlite
select
  g.name,
  ins.name as instance_name,
  ins.status as instance_status
from
  gcp_compute_instance_group as g,
  json_each(instances) as i,
  gcp_compute_instance as ins
where
  json_extract(i.value, '$.instance') = ins.self_link;
```

### Get network and subnetwork info of each instance group
Analyze the network configuration of each instance group to gain insights into the associated network and subnetwork details, such as the IP range, gateway address and location. This can be useful for assessing the network distribution of your resources.

```sql+postgres
select
  g.name as instance_group_name,
  n.name as network_name,
  s.name as subnetwork_name,
  s.ip_cidr_range,
  s.gateway_address,
  n.location
from
  gcp_compute_instance_group as g,
  gcp_compute_network as n,
  gcp_compute_subnetwork as s
where
  g.network = n.self_link
  and g.subnetwork = s.self_link;
```

```sql+sqlite
select
  g.name as instance_group_name,
  n.name as network_name,
  s.name as subnetwork_name,
  s.ip_cidr_range,
  s.gateway_address,
  n.location
from
  gcp_compute_instance_group as g,
  gcp_compute_network as n,
  gcp_compute_subnetwork as s
where
  g.network = n.self_link
  and g.subnetwork = s.self_link;
```

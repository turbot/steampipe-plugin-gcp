---
title: "Steampipe Table: gcp_compute_autoscaler - Query GCP Compute Autoscalers using SQL"
description: "Allows users to query GCP Compute Autoscalers, providing insights into autoscaling configurations and operational status."
folder: "Compute"
---

# Table: gcp_compute_autoscaler - Query GCP Compute Autoscalers using SQL

Google Cloud Platform's Compute Autoscaler is a service that automatically adjusts the number of virtual machine instances in a managed instance group based on the increase or decrease in load. It allows you to maintain the number of instances in your project so you can ensure that your application has the available resources to handle its work demand. This service is particularly useful for applications that experience variable levels of demand.

## Table Usage Guide

The `gcp_compute_autoscaler` table provides insights into the autoscaling configurations in Google Cloud Platform's Compute Engine. As a cloud engineer, explore autoscaler-specific details through this table, including the target that the autoscaler is configured to scale, the policy that defines the autoscaler behavior, and the operational status of the autoscaler. Utilize it to manage and optimize your resources, ensuring that your application has the necessary resources to handle its work demand.

## Examples

### Basic Info
Explore the basic details of your Google Cloud Platform's compute autoscalers, such as their names, statuses, and recommended sizes. This information can be useful in understanding your autoscaler setup and identifying areas for potential optimization.

```sql+postgres
select
  name,
  description,
  self_link,
  status,
  location,
  akas,
  recommended_size
from
  gcp_compute_autoscaler;
```

```sql+sqlite
select
  name,
  description,
  self_link,
  status,
  location,
  akas,
  recommended_size
from
  gcp_compute_autoscaler;
```

### Get auto scaling policy for autoscalers
Explore the configuration of auto scaling policies applied to autoscalers to better manage computational resources and optimize system performance.

```sql+postgres
select
  title,
  autoscaling_policy ->> 'mode' as mode,
  autoscaling_policy -> 'cpuUtilization' ->> 'predictiveMethod' as cpu_utilization_method,
  autoscaling_policy -> 'cpuUtilization' ->> 'utilizationTarget' as cpu_utilization_target,
  autoscaling_policy ->> 'maxNumReplicas' as max_replicas,
  autoscaling_policy ->> 'minNumReplicas' as min_replicas,
  autoscaling_policy ->> 'coolDownPeriodSec' as cool_down_period_sec
from
  gcp_compute_autoscaler;
```

```sql+sqlite
select
  title,
  json_extract(autoscaling_policy, '$.mode') as mode,
  json_extract(autoscaling_policy, '$.cpuUtilization.predictiveMethod') as cpu_utilization_method,
  json_extract(autoscaling_policy, '$.cpuUtilization.utilizationTarget') as cpu_utilization_target,
  json_extract(autoscaling_policy, '$.maxNumReplicas') as max_replicas,
  json_extract(autoscaling_policy, '$.minNumReplicas') as min_replicas,
  json_extract(autoscaling_policy, '$.coolDownPeriodSec') as cool_down_period_sec
from
  gcp_compute_autoscaler;
```

### Get autoscalers with configuration errors
Discover the autoscalers in your Google Cloud Platform that are encountering configuration errors. This can help in identifying and rectifying issues that could potentially hinder the automatic scaling of your resources.

```sql+postgres
select
  name,
  description,
  self_link,
  status,
  location,
  akas,
  recommended_size
from
  gcp_compute_autoscaler
where
  status = 'ERROR';
```

```sql+sqlite
select
  name,
  description,
  self_link,
  status,
  location,
  akas,
  recommended_size
from
  gcp_compute_autoscaler
where
  status = 'ERROR';
```

### Get instance groups having autoscaling enabled
Identify the groups of instances that have autoscaling enabled to ensure optimal resource allocation and cost efficiency. This is particularly useful for managing large-scale applications and maintaining performance during peak load times.

```sql+postgres
select
  a.title as autoscaler_name,
  g.name as instance_group_name,
  g.description as instance_group_description,
  g.size as instance_group_size
from
  gcp_compute_instance_group g,
  gcp_compute_autoscaler a
where
  g.name = split_part(a.target, 'instanceGroupManagers/', 2);
```

```sql+sqlite
select
  a.title as autoscaler_name,
  g.name as instance_group_name,
  g.description as instance_group_description,
  g.size as instance_group_size
from
  gcp_compute_instance_group g,
  gcp_compute_autoscaler a
where
  g.name = substr(a.target, instr(a.target, 'instanceGroupManagers/') + length('instanceGroupManagers/'));
```
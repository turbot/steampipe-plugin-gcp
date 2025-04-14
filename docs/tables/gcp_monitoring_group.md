---
title: "Steampipe Table: gcp_monitoring_group - Query GCP Monitoring Groups using SQL"
description: "Allows users to query Monitoring Groups in GCP, providing insights into various group configurations and their associated metrics."
folder: "Cloud Monitoring"
---

# Table: gcp_monitoring_group - Query GCP Monitoring Groups using SQL

A Monitoring Group in Google Cloud Platform (GCP) is a named set of Google Cloud resources identified by a filter. These groups provide a way to monitor and manage the combined behavior of a collection of related resources. They are useful for aggregate analysis and for building higher-level, more abstract collections.

## Table Usage Guide

The `gcp_monitoring_group` table provides insights into Monitoring Groups within Google Cloud Monitoring. As a DevOps engineer, explore group-specific details through this table, including resource type, group name, and associated metadata. Utilize it to uncover information about groups, such as those with specific configurations, the relationships between groups, and the verification of group metrics.

## Examples

### Filter info of each monitoring group
Explore which monitoring groups are in use in your Google Cloud Platform setup. This allows for better management of your resources by understanding the filters applied to each group.

```sql+postgres
select
  name,
  display_name,
  filter
from
  gcp_monitoring_group;
```

```sql+sqlite
select
  name,
  display_name,
  filter
from
  gcp_monitoring_group;
```

### List of cluster monitoring groups
Discover the segments that are grouped into clusters within the Google Cloud Platform's monitoring system, allowing you to better manage and organize your monitoring resources.

```sql+postgres
select
  name,
  display_name,
  is_cluster
from
  gcp_monitoring_group
where
  is_cluster;
```

```sql+sqlite
select
  name,
  display_name,
  is_cluster
from
  gcp_monitoring_group
where
  is_cluster = 1;
```
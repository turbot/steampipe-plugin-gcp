---
title: "Steampipe Table: gcp_vertex_ai_notebook_runtime_template - Query GCP Vertex AI Notebook Runtime Templates using SQL"
description: "Allows users to query GCP Vertex AI Notebook Runtime Templates, providing detailed information on configuration, machine specs, networking, and other related details."
folder: "Vertex AI"
---

# Table: gcp_vertex_ai_notebook_runtime_template - Query GCP Vertex AI Notebook Runtime Templates using SQL

Google Cloud Vertex AI Notebook Runtime Templates provide preconfigured environments for machine learning workloads. These templates allow you to define runtime settings such as machine type, networking, disk storage, and more. The `gcp_vertex_ai_notebook_runtime_template` table in Steampipe allows you to query detailed information about Notebook Runtime Templates, including network configurations, service account settings, and machine specifications.

## Table Usage Guide

The `gcp_vertex_ai_notebook_runtime_template` table helps cloud administrators, data scientists, and machine learning engineers gather insights into their Vertex AI Notebook Runtime Templates. You can query the configuration settings for machine specs, network access, disk storage, and more. This table is particularly useful for monitoring resource usage, managing templates, and ensuring that the runtime environment meets your workload needs.

## Examples

### Basic info
Retrieve basic information about the Notebook Runtime Templates, including their name, location, and creation time.

```sql+postgres
select
  name,
  display_name,
  location,
  create_time,
  is_default
from
  gcp_vertex_ai_notebook_runtime_template;
```

```sql+sqlite
select
  name,
  display_name,
  location,
  create_time,
  is_default
from
  gcp_vertex_ai_notebook_runtime_template;
```

### List templates by machine type
Identify templates that are using a specific machine type, such as "e2-standard-4".

```sql+postgres
select
  name,
  display_name,
  machine_spec ->> 'machine_type' as machine_type,
  location
from
  gcp_vertex_ai_notebook_runtime_template
where
  machine_spec ->> 'machine_type' = 'e2-standard-4';
```

```sql+sqlite
select
  name,
  display_name,
  json_extract(machine_spec, '$.machine_type') as machine_type,
  location
from
  gcp_vertex_ai_notebook_runtime_template
where
  json_extract(machine_spec, '$.machine_type') = 'e2-standard-4';
```

### List templates with internet access enabled
Retrieve templates that have internet access enabled in their network configuration.

```sql+postgres
select
  name,
  network_spec ->> 'enable_internet_access' as internet_access,
  network_spec ->> 'network' as network,
  location
from
  gcp_vertex_ai_notebook_runtime_template
where
  network_spec ->> 'enable_internet_access' = 'true';
```

```sql+sqlite
select
  name,
  json_extract(network_spec, '$.enable_internet_access') as internet_access,
  json_extract(network_spec, '$.network') as network,
  location
from
  gcp_vertex_ai_notebook_runtime_template
where
  json_extract(network_spec, '$.enable_internet_access') = 'true';
```

### List templates by disk size
Identify templates with specific persistent disk sizes, which is useful for understanding storage configurations.

```sql+postgres
select
  name,
  data_persistent_disk_spec ->> 'disk_size_gb' as disk_size_gb,
  location
from
  gcp_vertex_ai_notebook_runtime_template
where
  data_persistent_disk_spec ->> 'disk_size_gb' = '10';
```

```sql+sqlite
select
  name,
  json_extract(data_persistent_disk_spec, '$.disk_size_gb') as disk_size_gb,
  location
from
  gcp_vertex_ai_notebook_runtime_template
where
  json_extract(data_persistent_disk_spec, '$.disk_size_gb') = '10';
```

### List templates by tag
Retrieve templates associated with specific tags for organizational purposes.

```sql+postgres
select
  name,
  display_name,
  tags ->> 'foo' as tag_value,
  location
from
  gcp_vertex_ai_notebook_runtime_template
where
  tags ->> 'foo' = 'bar';
```

```sql+sqlite
select
  name,
  display_name,
  json_extract(tags, '$.foo') as tag_value,
  location
from
  gcp_vertex_ai_notebook_runtime_template
where
  json_extract(tags, '$.foo') = 'bar';
```
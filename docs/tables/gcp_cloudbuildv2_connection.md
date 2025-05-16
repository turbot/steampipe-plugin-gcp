---
title: "Steampipe Table: gcp_cloudbuildv2_connection - Query Cloud Build V2 Connections using SQL"
description: "Allows users to query Cloud Build V2 Connections in Google Cloud Platform (GCP), providing insights into connection configurations for various source code repositories and their integration states."
folder: "Cloud Build"
---

# Table: gcp_cloudbuildv2_connection - Query Cloud Build V2 Connections using SQL

Google Cloud Build V2 Connections enable integration between Cloud Build and various source code repositories such as GitHub, GitLab, and Bitbucket. These connections allow Cloud Build to automatically trigger builds when code changes are detected in the connected repositories. Cloud Build V2 Connections provide secure and managed access to your source code repositories.

## Table Usage Guide

The `gcp_cloudbuildv2_connection` table provides insights into Cloud Build V2 Connections within Google Cloud Platform. As a DevOps engineer or cloud administrator, explore connection-specific details through this table, including configuration settings, installation states, and associated metadata. Utilize it to manage and monitor your source code repository connections, verify connection states, and ensure proper integration with your CI/CD pipelines.

## Examples

### Basic info
Explore the basic details of your Cloud Build V2 connections, including their creation and update times, to understand when they were established and last modified.

```sql+postgres
select
  name,
  create_time,
  update_time,
  disabled,
  location
from
  gcp_cloudbuildv2_connection;
```

```sql+sqlite
select
  name,
  create_time,
  update_time,
  disabled,
  location
from
  gcp_cloudbuildv2_connection;
```

### List disabled connections
Identify instances where Cloud Build V2 connections are disabled, which can help in maintaining active and necessary connections while cleaning up unused ones.

```sql+postgres
select
  name,
  create_time,
  update_time,
  location
from
  gcp_cloudbuildv2_connection
where
  disabled;
```

```sql+sqlite
select
  name,
  create_time,
  update_time,
  location
from
  gcp_cloudbuildv2_connection
where
  disabled = 1;
```

### List GitHub connections with their configuration details
Analyze the settings of GitHub connections to understand their configuration details and ensure they are properly set up.

```sql
select
  name,
  github_config ->> 'host_uri' as host_uri,
  github_config ->> 'app_installation_id' as installation_id,
  github_config -> 'secrets' ->> 'oauth_token_secret_version' as oauth_token_secret,
  github_config -> 'secrets' ->> 'webhook_secret_secret_version' as webhook_secret
from
  gcp_cloudbuildv2_connection
where
  github_config is not null;
```

### Check connection installation states and action required
Analyze the installation states of your connections to identify any that require attention or have pending actions.

```sql+postgres
select
  name,
  installation_state ->> 'stage' as stage,
  installation_state ->> 'message' as message,
  installation_state ->> 'actionUri' as action_uri,
  location
from
  gcp_cloudbuildv2_connection
where
  installation_state ->> 'stage' != 'COMPLETE';
```

```sql+sqlite
select
  name,
  json_extract(installation_state, '$.stage') as stage,
  json_extract(installation_state, '$.message') as message,
  json_extract(installation_state, '$.actionUri') as action_uri,
  location
from
  gcp_cloudbuildv2_connection
where
  json_extract(installation_state, '$.stage') != 'COMPLETE';
```

### List connections by source control type with additional metadata
Analyze connections across different source control platforms and their associated metadata to get a comprehensive view of your build configurations.

```sql
select
  name,
  etag,
  case
    when github_config is not null then 'GitHub'
    when github_enterprise_config is not null then 'GitHub Enterprise'
    when gitlab_config is not null then 'GitLab'
    when bitbucket_data_center_config is not null then 'Bitbucket Data Center'
    when bitbucket_cloud_config is not null then 'Bitbucket Cloud'
  end as source_type,
  create_time,
  update_time,
  disabled,
  location,
  project
from
  gcp_cloudbuildv2_connection
order by
  source_type,
  name;
```
---
title: "Steampipe Table: gcp_app_engine_application - Query App Engine Application using SQL"
description: "Allows users to query App Engine Application in Google Cloud Platform (GCP), specifically the details about the application, including their name, location, service account, storage bucket, database type and serving_status."
folder: "App Engine"
---

# Table: gcp_app_engine_application - Query App Engine Application using SQL

Google Cloud Platform's (GCP) App Engine is a fully managed, serverless platform for developing and hosting web applications at scale. An App Engine application refers to the specific application you deploy on this platform.

## Table Usage Guide

The `gcp_app_engine_application` table provides insights into the App Engine abstracts away the infrastructure, allowing developers to focus on code. It manages the hardware and networking infrastructure required to run your code.

## Examples

### Basic info
Explore the basic details of your Google Cloud Platform's App Engine Application such as their names, storage bucket, database type, default hostname, and serving status. This information can help you manage and monitor your application more effectively.

```sql+postgres
select
  name,
  id,
  code_bucket,
  database_type,
  default_hostname,
  gcr_domain,
  serving_status
from
  gcp_app_engine_application;
```

```sql+sqlite
select
  name,
  id,
  code_bucket,
  database_type,
  default_hostname,
  gcr_domain,
  serving_status
from
  gcp_app_engine_application;
```

### Get feature setting details of an application
This is designed to retrieve specific configuration details from App Engine applications within a Google Cloud Platform (GCP) environment.

```sql+postgres
select
  name,
  id,
  location,
  feature_settings -> 'SplitHealthChecks' as split_health_checks,
  feature_settings -> 'UseContainerOptimizedOs' as use_container_optimized_os
from
  gcp_app_engine_application;
```

```sql+sqlite
select
  name,
  id,
  location,
  json_extract(feature_settings, '$.SplitHealthChecks') as split_health_checks,
  json_extract(feature_settings, '$.UseContainerOptimizedOs') as use_container_optimized_os
from
  gcp_app_engine_application;
```

### Get service account details for the application
Explore the details about the service account that has been associated with the application.

```sql+postgres
select
  a.name,
  a.service_account,
  s.email,
  s.disabled,
  s.oauth2_client_id,
  s.iam_policy
from
  gcp_app_engine_application as a,
  gcp_service_account as s
where
  s.name = a.service_account;
```

```sql+sqlite
select
  a.name,
  a.service_account,
  s.email,
  s.disabled,
  s.oauth2_client_id,
  s.iam_policy
from
  gcp_app_engine_application as a
join
  gcp_service_account as s ON s.name = a.service_account;
```
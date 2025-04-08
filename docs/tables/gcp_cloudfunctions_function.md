---
title: "Steampipe Table: gcp_cloudfunctions_function - Query Google Cloud Platform Cloud Functions using SQL"
description: "Allows users to query Cloud Functions in Google Cloud Platform, specifically the configuration, status, and metadata of each function, providing insights into the function's behavior and usage."
folder: "Cloud Functions"
---

# Table: gcp_cloudfunctions_function - Query Google Cloud Platform Cloud Functions using SQL

Google Cloud Functions is a lightweight, event-based, asynchronous compute solution that allows you to create small, single-purpose functions that respond to cloud events without the need to manage a server or a runtime environment. Events from Google Cloud Storage and Pub/Sub can trigger Cloud Functions asynchronously, or you can use HTTP invocation for synchronous execution. It allows developers to run their code in response to events without needing to provision or manage servers.

## Table Usage Guide

The `gcp_cloudfunctions_function` table provides insights into Cloud Functions within Google Cloud Platform. As a DevOps engineer, explore function-specific details through this table, including configuration, status, and associated metadata. Utilize it to uncover information about functions, such as their event triggers, resource usage, and execution environment.

## Examples

### Basic function info
Explore the operational status of various cloud functions to manage resources effectively. Analyze the settings to understand the runtime, available memory, and maximum instances for optimal performance.

```sql+postgres
select
  name,
  description,
  status,
  runtime,
  available_memory_mb,
  max_instances,
  ingress_settings,
  service_timeout
from
  gcp_cloudfunctions_function;
```

```sql+sqlite
select
  name,
  description,
  status,
  runtime,
  available_memory_mb,
  max_instances,
  ingress_settings,
  service_timeout
from
  gcp_cloudfunctions_function;
```

### Count of cloud functions by runtime engines
Analyze the distribution of cloud functions across different runtime engines, providing a useful overview to optimize resource allocation and understand usage patterns.
```sql+postgres
select
  runtime,
  count(*)
from
  gcp_cloudfunctions_function
group by
  runtime;
```

```sql+sqlite
select
  runtime,
  count(*)
from
  gcp_cloudfunctions_function
group by
  runtime;
```

### Cloud functions service account info
Explore which cloud functions are linked to specific service accounts. This can help manage and secure access to resources, by ensuring only authorized accounts are connected to specific functions.

```sql+postgres
select
  f.name as function_name,
  f.service_account_email as service_account_email,
  a.display_name as service_account_display_name
from
  gcp_cloudfunctions_function as f,
  gcp_service_account as a
where
  f.service_account_email = a.email;
```

```sql+sqlite
select
  f.name as function_name,
  f.service_account_email as service_account_email,
  a.display_name as service_account_display_name
from
  gcp_cloudfunctions_function as f,
  gcp_service_account as a
where
  f.service_account_email = a.email;
```

### Cloud functions service account info, including roles assigned in the project IAM policy
Determine the roles assigned to various service accounts within your project's IAM policy, particularly those associated with cloud functions. This can help maintain security by ensuring only necessary permissions are granted.

```sql+postgres
select
  f.name as function_name,
  f.service_account_email as service_account_email,
  a.display_name as service_account_display_name,
  b ->> 'role' as role_name
from
  gcp_cloudfunctions_function as f,
  gcp_service_account as a,
  gcp_iam_policy as p,
  jsonb_array_elements(bindings) as b,
  jsonb_array_elements_text(b -> 'members') as m
where
  f.service_account_email = a.email
  and m = ( 'serviceAccount:' || f.service_account_email);
```

```sql+sqlite
select
  f.name as function_name,
  f.service_account_email as service_account_email,
  a.display_name as service_account_display_name,
  json_extract(b.value, '$.role') as role_name
from
  gcp_cloudfunctions_function as f,
  gcp_service_account as a,
  gcp_iam_policy as p,
  json_each(bindings) as b,
  json_each(json_extract(b.value, '$.members')) as m
where
  f.service_account_email = a.email
  and m.value = ( 'serviceAccount:' || f.service_account_email);
```

### View the resource-level IAM policy on cloud functions
Explore the access control measures applied to your cloud functions. This query is useful to understand the security configuration and permissions associated with each function, helping to maintain robust access management.
```sql+postgres
select
  name,
  jsonb_pretty(iam_policy)
from
  gcp_cloudfunctions_function;
```

```sql+sqlite
select
  name,
  iam_policy
from
  gcp_cloudfunctions_function;
```

### Find members assigned in resource-level IAM policy on cloud functions that are not in your email domain
Explore which members are assigned in resource-level IAM policy on cloud functions that are not within your email domain. This is useful to identify potential security risks by detecting unauthorized users who might have access to your cloud functions.

```sql+postgres
select
  name,
  b ->> 'role' as role_name,
  m as member
from
  gcp_cloudfunctions_function,
  jsonb_array_elements(iam_policy -> 'bindings') as b,
  jsonb_array_elements_text(b -> 'members') as m
where
  m not like '%@turbot.com';
```

```sql+sqlite
select
  name,
  json_extract(b.value, '$.role') as role_name,
  m.value as member
from
  gcp_cloudfunctions_function,
  json_each(iam_policy, '$.bindings') as b,
  json_each(b.value, '$.members') as m
where
  m.value not like '%@turbot.com';
```
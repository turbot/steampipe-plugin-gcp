# Table:  gcp_cloudfunctions_function

Google Cloud Functions is a serverless execution environment for building and connecting cloud services. With Cloud Functions you write simple, single-purpose functions that are attached to events emitted from your cloud infrastructure and services. 

## Examples


### Basic function info

```sql
select
  name,
  description,
  status,
  runtime,
  available_memory_mb,
  max_instances,
  ingress_settings,
  timeout
from
  gcp_cloudfunctions_function;
```


### Count of cloud functions by runtime engines
```sql
select
  runtime,
  count(*)
from
  gcp_cloudfunctions_function
group by
  runtime;
```


### Cloud functions service account info
```sql
select
  f.name as function_name,
  f.service_account_email as service_account_email,
  a.display_name as service_account_display_name
from
  gcp_cloudfunctions_function as f,
  gcp_service_account as a
where 
  f.service_account_email = a.email
```


### Cloud functions service account info, including roles assigned in the project IAM policy
```sql
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
  and m = ( 'serviceAccount:' || f.service_account_email)
```


### View the resource-level IAM policy on cloud functions
```sql
select
  name,
  jsonb_pretty(iam_policy)
from
  gcp_cloudfunctions_function;
```

### Find members assigned in resource-level IAM policy on cloud functions that are not in your email domain

```sql
select
  name,
  b ->> 'role' as role_name,
  m as member
from
  gcp_cloudfunctions_function,
  jsonb_array_elements(iam_policy -> 'bindings') as b,
  jsonb_array_elements_text(b -> 'members') as m
where
  m not like '%@turbot.com'
```

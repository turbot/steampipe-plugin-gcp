---
title: "Steampipe Table: gcp_apikeys_key - Query GCP API Keys using SQL"
description: "Allows users to query API Keys in GCP, specifically the details of each API key and its associated metadata, providing insights into key usage and permissions."
folder: "API Keys"
---

# Table: gcp_apikeys_key - Query GCP API Keys using SQL

Google Cloud Platform (GCP) API Keys are unique identifiers used to authenticate users, applications, or devices to your APIs. They are used to track and control how the API is used, for example, to prevent malicious use or abuse of your APIs. API keys are project-centric, meaning they are created, managed, and used by APIs within a specific project.

## Table Usage Guide

The `gcp_apikeys_key` table provides insights into API keys within Google Cloud Platform (GCP). As a developer or security analyst, explore key-specific details through this table, including permissions, creation time, and associated metadata. Utilize it to uncover information about keys, such as those with specific permissions, the usage of keys, and the verification of key restrictions.

## Examples

### Basic info
Explore which API keys in your Google Cloud Platform have restrictions. This allows you to determine the state of each key, providing insights into their creation time and the level of access they provide.

```sql+postgres
select
  uid,
  display_name,
  create_time,
  case when restrictions is null then 'Unrestricted' else 'Restricted' end as state
from
  gcp_apikeys_key;
```

```sql+sqlite
select
  uid,
  display_name,
  create_time,
  case when restrictions is null then 'Unrestricted' else 'Restricted' end as state
from
  gcp_apikeys_key;
```


### List all unrestricted keys
Explore which API keys in your Google Cloud Platform account have no set restrictions, allowing you to identify potential security risks. This can be useful in assessing the elements within your environment that may be open to misuse or unauthorized access.

```sql+postgres
select
  uid,
  display_name,
  create_time,
  case when restrictions is null then 'Unrestricted' else 'Restricted' end as state
from
  gcp_apikeys_key
where
  restrictions is null;
```

```sql+sqlite
select
  uid,
  display_name,
  create_time,
  case when restrictions is null then 'Unrestricted' else 'Restricted' end as state
from
  gcp_apikeys_key
where
  restrictions is null;
```

### Get api service restrictions associated with each key
Determine the restrictions linked with each API key to understand the allowed services. This can help manage access and maintain security by identifying which services are accessible with each key.

```sql+postgres
select
  uid,
  display_name,
  a ->> 'service' as allowed_service
from
  gcp_apikeys_key,
  jsonb_array_elements(restrictions -> 'apiTargets') as a
where
  restrictions is not null;
```

```sql+sqlite
select
  uid,
  display_name,
  json_extract(a.value, '$.service') as allowed_service
from
  gcp_apikeys_key,
  json_each(restrictions, '$.apiTargets') as a
where
  restrictions is not null;
```

### Get website restrictions associated with each key
Determine the areas in which each key has associated website restrictions. This query is useful in understanding the limitations set on each key, providing insights into potential access or usage constraints.

```sql+postgres
select
  uid,
  display_name,
  a as allowed_website
from
  gcp_apikeys_key,
  jsonb_array_elements_text(restrictions -> 'browserKeyRestrictions' -> 'allowedReferrers') as a
where
  restrictions is not null;
```

```sql+sqlite
select
  uid,
  display_name,
  json_extract(restrictions, '$.browserKeyRestrictions.allowedReferrers') as allowed_website
from
  gcp_apikeys_key
where
  json_type(restrictions, '$.browserKeyRestrictions.allowedIps') = 'array';
```

### Get ip restrictions associated with each key
Explore which API keys have associated IP restrictions in your Google Cloud Platform. This can help in identifying potential security risks and ensuring that only authorized IPs have access to your keys.

```sql+postgres
select
  uid,
  display_name,
  a as allowed_ip
from
  gcp_apikeys_key,
  jsonb_array_elements_text(restrictions -> 'serverKeyRestrictions' -> 'allowedIps') as a
where
  restrictions is not null;
```

```sql+sqlite
select
  uid,
  display_name,
  json_extract(restrictions, '$.serverKeyRestrictions.allowedIps') as allowed_ip
from
  gcp_apikeys_key
where
  json_type(restrictions, '$.serverKeyRestrictions.allowedIps') = 'array';
```

### Get iOS app restrictions associated with each key
Discover the segments that indicate the restrictions placed on each iOS application associated with a specific key. This can help in managing app permissions and ensuring the security of your digital assets.

```sql+postgres
select
  uid,
  display_name,
  a as allowed_ios_bundle_id
from
  gcp_apikeys_key,
  jsonb_array_elements_text(restrictions -> 'iosKeyRestrictions' -> 'allowedBundleIds') as a
where
  restrictions is not null;
```

```sql+sqlite
select
  uid,
  display_name,
  a.value as allowed_ios_bundle_id
from
  gcp_apikeys_key,
  json_each(restrictions, '$.iosKeyRestrictions.allowedBundleIds') as a
where
  restrictions is not null;
```

### Get android app restrictions associated with each key
Identify the restrictions associated with each Android application in your Google Cloud Platform. This can help in managing and controlling access to your applications, thus enhancing security.

```sql+postgres
select
  uid,
  display_name,
  a as allowed_android_apps
from
  gcp_apikeys_key,
  jsonb_array_elements(restrictions -> 'androidKeyRestrictions' -> 'allowedApplications') as a
where
  restrictions is not null;
```

```sql+sqlite
select
  uid,
  display_name,
  json_extract(restrictions, '$.androidKeyRestrictions.allowedApplications') as allowed_ip
from
  gcp_apikeys_key
where
  json_type(restrictions, '$.androidKeyRestrictions.allowedApplications') = 'array';
```
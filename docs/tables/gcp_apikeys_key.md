# Table: gcp_apikeys_key

An API key is a simple encrypted string that you can use when calling Google Cloud APIs. A typical use of an API key is to pass the key into a REST API call as a query parameter with the following format:

### Basic info

```sql
select
  uid,
  display_name,
  create_time,
  case when restrictions is null then 'Unrestricted' else 'Restricted' end as state
from
  gcp_apikeys_key;
```


### List all unrestricted keys

```sql
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

```sql
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

### Get website restrictions associated with each key

```sql
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

### Get ip restrictions associated with each key

```sql
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

### Get iOS app restrictions associated with each key

```sql
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

### Get android app restrictions associated with each key

```sql
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
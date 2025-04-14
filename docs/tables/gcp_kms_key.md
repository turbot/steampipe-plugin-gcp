---
title: "Steampipe Table: gcp_kms_key - Query Google Cloud KMS Keys using SQL"
description: "Allows users to query Google Cloud KMS Keys, providing insights into key management and encryption details."
folder: "KMS"
---

# Table: gcp_kms_key - Query Google Cloud KMS Keys using SQL

Google Cloud Key Management Service (KMS) is a cloud service for managing encryption keys on Google Cloud. This service allows you to generate, use, rotate, and destroy AES256, RSA 2048, RSA 3072, RSA 4096, EC P256, and EC P384 cryptographic keys. KMS is integrated with Cloud IAM and Cloud Audit Logging so that you can manage permissions on individual keys and monitor how these are used.

## Table Usage Guide

The `gcp_kms_key` table provides insights into the cryptographic keys managed by Google Cloud KMS. As a security engineer, you can explore key-specific details through this table, including key versions, key state, and associated metadata. Utilize it to uncover information about key usage, rotation schedule, and the verification of key permissions.

## Examples

### Basic info
Explore which cryptographic keys in the Google Cloud Platform have been created and their rotation periods. This can be useful in managing the lifecycle of your keys and ensuring regular key rotation for enhanced security.

```sql+postgres
select
  name,
  create_time,
  rotation_period
from
  gcp_kms_key;
```

```sql+sqlite
select
  name,
  create_time,
  rotation_period
from
  gcp_kms_key;
```

### List keys older than 30 days
Explore which security keys have been in use for more than a month. This can help in maintaining security standards by regularly updating the keys.

```sql+postgres
select
  name,
  create_time,
  rotation_period
from
  gcp_kms_key
where
  create_time <= (current_date - interval '30' day)
order by
  create_time;
```

```sql+sqlite
select
  name,
  create_time,
  rotation_period
from
  gcp_kms_key
where
  date(create_time) <= date('now','-30 day')
order by
  create_time;
```

### List keys with rotation period greater than 90 days (7776000 seconds)
Determine the areas in which encryption keys have a rotation period exceeding 90 days, a parameter that may be relevant for assessing the security measures in place within your GCP environment.

```sql+postgres
select
  name,
  create_time,
  rotation_period
from
  gcp_kms_key
where
  split_part(rotation_period, 's', 1) :: int > 7776000;
```

```sql+sqlite
select
  name,
  create_time,
  rotation_period
from
  gcp_kms_key
where
  cast(substr(rotation_period, 1, instr(rotation_period, 's') - 1) as integer) > 7776000;
```

### List publicly accessible keys
The query helps identify any security risks by pinpointing instances where encryption keys are publicly accessible in your Google Cloud Platform. This can assist in maintaining data confidentiality and preventing unauthorized access.

```sql+postgres
select
  distinct name,
  key_ring_name,
  location
from
  gcp_kms_key,
  jsonb_array_elements(iam_policy -> 'bindings') as b
where
  b -> 'members' ?| array['allAuthenticatedUsers', 'allUsers'];
```

```sql+sqlite
select distinct
  k.name,
  k.key_ring_name,
  k.location
from
  gcp_kms_key k,
  json_each(k.iam_policy, '$.bindings') as b
where
  json_extract(b.value, '$.members') like '%allAuthenticatedUsers%' OR
  json_extract(b.value, '$.members') like '%allUsers%';
```
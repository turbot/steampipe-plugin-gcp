---
title: "Steampipe Table: gcp_kms_key_version - Query Google Cloud KMS Key Versions using SQL"
description: "Allows users to query Google Cloud KMS Key Versions, specifically the key version details, providing insights into key management and security compliance."
folder: "KMS"
---

# Table: gcp_kms_key_version - Query Google Cloud KMS Key Versions using SQL

Google Cloud Key Management Service (KMS) is a cloud service for managing cryptographic keys for your cloud services the same way you do on-premises. It provides the capability to create, import, manage, rotate, and destroy AES256, RSA 2048, RSA 3072, RSA 4096, EC P256, and EC P384 cryptographic keys. KMS is integrated with Cloud IAM and Cloud Audit Logging so that you can manage permissions on individual keys and monitor how these are used.

## Table Usage Guide

The `gcp_kms_key_version` table provides insights into key versions within Google Cloud Key Management Service (KMS). As a security or compliance professional, explore key version-specific details through this table, including key material, state, and associated metadata. Utilize it to uncover information about key versions, such as those in use, the cryptographic configuration of each key version, and the verification of their lifecycle state.

## Examples

### Basic info
Explore the status of encryption keys in Google Cloud Platform's Key Management Service, excluding those that have been destroyed. This can be useful in identifying active keys and ensuring proper key management.

```sql+postgres
select
  key_name,
  crypto_key_version,
  title,
  state
from
  gcp_kms_key_version
where
  state <> 'DESTROYED';
```

```sql+sqlite
select
  key_name,
  crypto_key_version,
  title,
  state
from
  gcp_kms_key_version
where
  state <> 'DESTROYED';
```

### List key versions older than 30 days
Explore which key versions in Google Cloud's Key Management Service are older than 30 days and have not been destroyed. This can help identify outdated keys that might need updating or deletion for security purposes.

```sql+postgres
select
  key_name,
  create_time,
  crypto_key_version,
  state
from
  gcp_kms_key_version
where
  create_time <= (current_date - interval '30' day) and
  state <> 'DESTROYED'
order by
  create_time;
```

```sql+sqlite
select
  key_name,
  create_time,
  crypto_key_version,
  state
from
  gcp_kms_key_version
where
  date(create_time) <= date('now', '-30 days')
  and state <> 'DESTROYED'
order by
  create_time;
```

### List key versions using google symmetric encryption algorithm
Explore which encryption keys are using the Google Symmetric Encryption algorithm. This can help you assess the security measures in place and ensure that they are up to date.

```sql+postgres
select
  key_name,
  create_time,
  crypto_key_version,
  algorithm
from
  gcp_kms_key_version
where
  algorithm like 'GOOGLE_SYMMETRIC_ENCRYPTION'
order by
  create_time;
```

```sql+sqlite
select
  key_name,
  create_time,
  crypto_key_version,
  algorithm
from
  gcp_kms_key_version
where
  algorithm like 'GOOGLE_SYMMETRIC_ENCRYPTION'
order by
  create_time;
```

### List disabled keys
Analyze the settings to understand which keys have been disabled in the GCP Key Management System. This can be useful for identifying potential security risks and ensuring proper key management.

```sql+postgres
select
  key_name,
  max(crypto_key_version) crypto_key_version,
  state
from
  gcp_kms_key_version
where
  state like 'DISABLED'
group by
  key_name,
  state;
```

```sql+sqlite
select
  key_name,
  max(crypto_key_version) as crypto_key_version,
  state
from
  gcp_kms_key_version
where
  state like 'DISABLED'
group by
  key_name,
  state;
```
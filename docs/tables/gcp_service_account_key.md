---
title: "Steampipe Table: gcp_service_account_key - Query Google Cloud Platform Service Account Keys using SQL"
description: "Allows users to query Service Account Keys in Google Cloud Platform, providing details about the keys associated with service accounts."
folder: "Service Account"
---

# Table: gcp_service_account_key - Query Google Cloud Platform Service Account Keys using SQL

A Service Account Key in Google Cloud Platform is a cryptographic key associated with a service account that can be used to authenticate as the service account. Service Account Keys are used to sign tokens for service accounts. They are essential for server-to-server interactions that are not tied to a user identity.

## Table Usage Guide

The `gcp_service_account_key` table provides insights into Service Account Keys within Google Cloud Platform. As a security engineer, explore key-specific details through this table, including the associated service account, key algorithm, and key origin. Utilize it to understand the distribution of keys, their validity, and their associated service accounts for better management and security.

## Examples

### List of service accounts using user managed keys
Identify the service accounts that utilize user-managed keys. This is useful to gain insights into potential security risks, as these keys are not automatically rotated and require manual management.

```sql+postgres
select
  service_account_name as service_account,
  title,
  key_type
from
  gcp_service_account_key
where
  key_type = 'USER_MANAGED';
```

```sql+sqlite
select
  service_account_name as service_account,
  title,
  key_type
from
  gcp_service_account_key
where
  key_type = 'USER_MANAGED';
```

### Validity time for the service account keys
Assess the elements within your Google Cloud Platform by identifying the validity period of your service account keys. This allows you to manage access and security by knowing when these keys are active.

```sql+postgres
select
  title,
  service_account_name as service_account,
  valid_after_time,
  valid_before_time
from
  gcp_service_account_key;
```

```sql+sqlite
select
  title,
  service_account_name as service_account,
  valid_after_time,
  valid_before_time
from
  gcp_service_account_key;
```

### Get public key data for a service account key
Explore the public key data associated with a specific service account key, allowing you to gain insights into the key type, origin, and format. This can be useful for verifying the key's authenticity and ensuring its proper configuration.

```sql+postgres
select
  name,
  key_type,
  key_origin,
  public_key_data_raw,
  public_key_data_pem
from
  gcp_service_account_key
where
  service_account_name = 'test@myproject.iam.gserviceaccount.com';
```

```sql+sqlite
select
  name,
  key_type,
  key_origin,
  public_key_data_raw,
  public_key_data_pem
from
  gcp_service_account_key
where
  service_account_name = 'test@myproject.iam.gserviceaccount.com';
```
---
title: "Steampipe Table: gcp_kms_key_ring - Query Google Cloud Key Management Service Key Rings using SQL"
description: "Allows users to query Google Cloud KMS Key Rings, providing detailed information about each key ring including its name, creation time, and location."
folder: "KMS"
---

# Table: gcp_kms_key_ring - Query Google Cloud Key Management Service Key Rings using SQL

Google Cloud Key Management Service (KMS) is a cloud service for managing cryptographic keys for your cloud services. Key Rings are used to group keys together for easier management. Each Key Ring belongs to a specific Google Cloud Project and resides in a specific location.

## Table Usage Guide

The `gcp_kms_key_ring` table provides insights into Key Rings within Google Cloud Key Management Service (KMS). As a security engineer, you can explore key ring-specific details through this table, including their names, creation times, and locations. Utilize it to manage and review the cryptographic keys for your cloud services.

## Examples

### Basic info
Explore which key rings have been created within Google Cloud's Key Management Service. This can help monitor the timeline of key ring creation for better security and resource management.

```sql+postgres
select
  name,
  create_time
from
  gcp_kms_key_ring;
```

```sql+sqlite
select
  name,
  create_time
from
  gcp_kms_key_ring;
```

### List key rings older than 30 days
Discover the segments that have key rings older than 30 days. This is useful for identifying and managing outdated or potentially unused key rings in your Google Cloud Platform.

```sql+postgres
select
  name,
  create_time
from
  gcp_kms_key_ring
where
  create_time <= (current_date - interval '30' day)
order by
  create_time;
```

```sql+sqlite
select
  name,
  create_time
from
  gcp_kms_key_ring
where
  date(create_time) <= date('now', '-30 days')
order by
  create_time;
```
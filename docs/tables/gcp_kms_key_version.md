# Table: gcp_kms_key_version

In Cloud KMS, the cryptographic key material that you use to encrypt, decrypt, sign, and verify data is stored in a key version. A key has zero or more key versions. When you rotate a key, you create a new key version.

## Examples

### Basic info

```sql
select
  name,
  crypto_key_version,
  state
from
  gcp_kms_key_version
where
  state <> 'DESTROYED';
```

### List key versions older than 30 days

```sql
select
  name,
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

### List key versions using google symmetric encryption algorithm

```sql
select
  name,
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

```sql
select 
  name,
  max(crypto_key_version) crypto_key_version, 
  state 
from 
  gcp_kms_key_version 
where 
  state like 'DISABLED' 
group by 
  name,
  state;
```
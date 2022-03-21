# Table: gcp_kms_key_version

Each version of a key contains key material used for encryption or signing. A key's version is represented by an integer, starting at 1.

## Examples

### Basic info

```sql
select
  name,
  crypto_key_version,
  state
from
  gcp_kms_key_version;
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
  create_time <= (current_date - interval '30' day)
order by
  create_time;
```

### List key versions that are enabled

```sql
select
  name,
  create_time,
  crypto_key_version,
  state
from
  gcp_kms_key_version
where
  state like 'ENABLED'
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
### List latest key versions that are enabled for crypto keys

```sql
select 
  name,
  max(crypto_key_version) crypto_key_version, 
  state 
from 
  gcp_kms_key_version 
where 
  state like 'ENABLED' 
group by 
  name,
  state;
```
# Table: gcp_service_account_key

Service Account Keys are public/private RSA key pairs which are used to authenticate to Google Cloud APIs.

## Examples

### List of service accounts using user managed keys

```sql
select
  split_part(name, '/', 4) as service_accounts,
  title,
  key_type,
from
  gcp_service_account_key
where
  key_type = 'USER_MANAGED';
```


### Validity time for the service account keys

```sql
select
  title,
  split_part(name, '/', 4) as service_accounts,
  valid_after_time,
  valid_before_time
from
  gcp_service_account_key;
```

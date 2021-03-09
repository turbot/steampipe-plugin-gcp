# Table: gcp_service_account_key

Service Account Keys are public/private RSA key pairs which are used to authenticate to Google Cloud APIs.

## Examples

### List of service accounts using user managed keys

```sql
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

```sql
select
  title,
  service_account_name as service_account,
  valid_after_time,
  valid_before_time
from
  gcp_service_account_key;
```

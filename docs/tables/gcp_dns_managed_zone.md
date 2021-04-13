# Table: gcp_dns_managed_zone

A DNS zone is used to host the DNS records for a particular domain. To start hosting your domain in Azure DNS, you need to create a DNS zone for that domain name. Each DNS record for your domain is then created inside this DNS zone.

## Examples

### Basic info

```sql
select
  name,
  id,
  dns_name,
  creation_time,
  visibility
from
  gcp_dns_managed_zone;
```

### List public zones with DNSSEC disabled

```sql
select
  name,
  id,
  dns_name,
  dnssec_config_state,
  visibility
from
  gcp_dns_managed_zone
where 
  visibility = 'public'
  and 
  (
    dnssec_config_state is null
    or dnssec_config_state = 'off'
  );
```

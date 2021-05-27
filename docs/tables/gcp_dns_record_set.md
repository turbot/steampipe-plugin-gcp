# Table: gcp_dns_record_set

A record set (also known as a resource record set) is the collection of DNS records in a zone that have the same name and are of the same type. Most record sets contain a single record.

## Examples

### Basic info

```sql
select
  name, 
  managed_zone_name, 
  type, 
  kind, 
  routing_policy,
  rrdatas,
  signature_rrdatas,
  ttl
from
  gcp_dns_record_set;
```

### List record sets of type 'CNAME'

```sql
select
  name, 
  managed_zone_name, 
  type, 
  ttl
from
  gcp_dns_record_set
where 
 type = 'CNAME';
```
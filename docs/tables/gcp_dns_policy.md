# Table: gcp_dns_managed_zone

You can configure one DNS server policy for each Virtual Private Cloud (VPC) network. The policy can specify inbound DNS forwarding, outbound DNS forwarding, or both. In this section, inbound server policy refers to a policy that permits inbound DNS forwarding. Outbound server policy refers to one possible method for implementing outbound DNS forwarding. It is possible for a policy to be both an inbound server policy and an outbound server policy if it implements the features of both.

## Examples

### Basic info

```sql
select
  name, 
  title, 
  id, 
  kind, 
  description, 
  enable_inbound_forwarding, 
  enable_logging, 
  target_name_servers 
from 
  gcp_dns_policy 
```

### List dns policies with EnableLogging enabled

```sql
select 
  name, 
  title, 
  id, 
  kind, 
  description, 
  enable_inbound_forwarding, 
  enable_logging 
from 
  gcp_dns_policy 
where 
  enable_logging 
```

# Table: gcp_dns_policy

You can configure one DNS server policy for each Virtual Private Cloud (VPC) network. The policy can specify inbound DNS forwarding, outbound DNS forwarding, or both.
Inbound server policy refers to a policy that permits inbound DNS forwarding, and outbound server policy refers to one possible method for implementing outbound DNS forwarding. It is possible for a policy to be both an inbound server policy and an outbound server policy if it implements the features of both.

## Examples

### Basic info

```sql
select
  name,
  id,
  kind,
  enable_inbound_forwarding,
  enable_logging,
  target_name_servers
from
  gcp_dns_policy;
```

### List dns policies with logging disabled

```sql
select
  name,
  id,
  enable_logging
from
  gcp_dns_policy
where
  not enable_logging;
```

### List dns policies not associated with any network

```sql
select
  name,
  id,
  networks
from
  gcp_dns_policy
where
  networks = '[]';
```

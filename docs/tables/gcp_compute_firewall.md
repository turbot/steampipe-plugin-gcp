# Table: gcp_compute_firewall

VPC firewall rules allows or denies connections to or from your virtual machine (VM) instances based on a specified configuration. Enabled VPC firewall rules are always enforced, protecting instances regardless of their configuration and operating system, even if they have not started up.

### Firewall rules basic info

```sql
select
  name,
  id,
  description,
  direction
from
  gcp_compute_firewall;
```


### List of rules which are applied to TCP protocol

```sql
select
  name,
  id,
  p ->> 'IPProtocol' as ip_protocol,
  p ->> 'ports' as ports
from
  gcp_compute_firewall,
  jsonb_array_elements(allowed) as p
where
  p ->> 'IPProtocol' = 'tcp';
```


### List of disabled rules

```sql
select
  name,
  id,
  description,
  disabled
from
  gcp_compute_firewall
where
  disabled
```


### List of Egress rules

```sql
select
  name,
  id,
  direction,
  allowed,
  denied
from
  gcp_compute_firewall
where
  direction = 'EGRESS';
```

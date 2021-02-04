# Table: gcp_compute_subnetwork

A subnetwork (also known as a subnet) is a logical partition of a Virtual Private Cloud network with one primary IP range and zero or more secondary IP ranges.

## Examples

### Subnetwork basic info

```sql
select
  name,
  id,
  network_name
from
  gcp_compute_subnetwork;
```

### List subnetworks having VPC flow logging set to false

```sql
select
  name,
  id,
  enable_flow_logs
from
  gcp_compute_subnetwork
where
  enable_flow_logs = false;
```

### List IAM policy attached with subnetworks

```sql
select
  name,
  iam_policy
from
  gcp_compute_subnetwork;
```

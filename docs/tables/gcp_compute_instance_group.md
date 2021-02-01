# Table:  gcp_compute_instance_group

An instance group is a collection of virtual machine (VM) instances that you can manage as a single entity.

## Examples

### Basic info

```sql
select
  name,
  kind,
  size,
  self_link
from
  gcp_compute_instance_group;
```

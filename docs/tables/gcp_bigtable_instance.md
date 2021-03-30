# Table: gcp_bigtable_instance

A Cloud Bigtable instance is a container for your data. Instances have one or more clusters, located in different zones. Each cluster has at least 1 node.

## Examples

### Basic info

```sql
select
  name,
  instance_type,
  state,
  location
from
  gcp_bigtable_instance;
```

### Get members and their associated IAM roles for each instance

```sql
select
  name,
  location,
  jsonb_array_elements_text(p -> 'members') as member,
  p ->> 'role' as role
from
  gcp_bigtable_instance,
  jsonb_array_elements(iam_policy -> 'bindings') as p;
```

### List instances whose members have Bigtable admin access

```sql
select
  name,
  instance_type,
  jsonb_array_elements_text(i -> 'members') as members,
  i ->> 'role' as role
from
  gcp_bigtable_instance,
  jsonb_array_elements(iam_policy -> 'bindings') as i
where
  i ->> 'role' like '%bigtable.admin';
```

### Count the number of instances per instance type

```sql
select
  instance_type,
  count(name)
from
  gcp_bigtable_instance
group by
  instance_type;
```

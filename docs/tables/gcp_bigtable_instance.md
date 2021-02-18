# Table: gcp_bigtable_instance

A Cloud Bigtable instance is a container for your data. Instances have one or more clusters, located in different zones. Each cluster has at least 1 node.

## Examples

### Basic info

```sql
select
  name,
  instance_type,
  location
from
  gcp_bigtable_instance;
```


### List of bigtable instances where members have IAM security admin access assigned in a resource policy

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
  i ->> 'role' like '%securityAdmin%';
```


### Get the count of instances which are not for production

```sql
select
  instance_type,
  count(name)
from
  gcp_bigtable_instance
where
  instance_type <> 'PRODUCTION'
group by
  instance_type;
```

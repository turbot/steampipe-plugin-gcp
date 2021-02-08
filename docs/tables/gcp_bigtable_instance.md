# Table: gcp_bigtable_instance

A Cloud Bigtable instance is a container for your data. Instances have one or more clusters, located in different zones. Each cluster has at least 1 node.

## Examples

### Basic info

```sql
select
  name,
  instance_type
from
  gcp_bigtable_instance;
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
group by instance_type;
```

# Table:  gcp_monitoring_group

Groups provide a mechanism for alerting on the behavior of a set of resources, rather than on individual resources

## Examples

### Filter info of each monitoring group

```sql
select
  name,
  display_name,
  filter
from
  gcp_monitoring_group;
```


### List of cluster monitoring groups

```sql
select
  name,
  display_name,
  is_cluster
from
  gcp_monitoring_group
where
  is_cluster;
```

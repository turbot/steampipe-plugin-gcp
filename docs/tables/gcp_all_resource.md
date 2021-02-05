# Table: gcp_all_resource

Lists all Cloud resources within the specified scope (i.e. Project).

## Examples

### Get the count of resource of specific resource type

```sql
select
  type,
  count(name)
from
  gcp_all_resource
group by type;
```

### List of all resource in current GCP project

```sql
select
  title,
  type,
  project
from
  gcp_all_resource;
```

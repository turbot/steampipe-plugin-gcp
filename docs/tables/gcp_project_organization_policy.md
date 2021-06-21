# Table: gcp_project_organization_policy

The Organization Policy Service gives you centralized and programmatic control over your organization's cloud resources.

## Examples

### Basic info

```sql
select
  id,
  version,
  update_time
from
  gcp_project_organization_policy;
```

### Get organization policy constraints for each policy

```sql
select
  id,
  version,
  list_policy ->> 'allValues' as policy_value
from
  gcp_project_organization_policy;
```

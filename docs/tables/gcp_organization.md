# Table: gcp_organization

GCP Organization resource is the root node of the GCP (Google Cloud Platform) resource hierarchy and all resources that belong to an organization are located under the organization node.

## Examples

### Basic info

```sql
select
  display_name,
  organization_id,
  lifecycle_state,
  creation_time
from
  gcp_organization;
```

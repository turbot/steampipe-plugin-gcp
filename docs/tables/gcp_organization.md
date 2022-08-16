# Table: gcp_organization

A GCP organization is the root node of the GCP (Google Cloud Platform) resource hierarchy and all resources that belong to an organization are located under the organization node.

**_Please note_**: This table requires the `resourcemanager.organizations.get` permission to retrieve organization details.

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

### Get essential contacts for organizations

```sql
Select
  organization_id,
  jsonb_pretty(essential_contacts) as essential_contacts
from
  gcp_organization;
```

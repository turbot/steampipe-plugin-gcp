# Table: gcp_project

A project organizes all your Google Cloud resources. A project consists of a set of users, a set of APIs and billing, authentication and monitoring settings for those APIs.

## Examples

### Basic info

```sql
select
  name,
  project_id,
  project_number,
  lifecycle_state,
  create_time
from
  gcp_project;
```

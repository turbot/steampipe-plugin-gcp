# Table: gcp_project_service

Services which are enabled or disabled in the project.

### Basic info

```sql
select
  *
from
  gcp_project_service;
```


### List of services which are enabled

```sql
select
  name,
  state
from
  gcp_project_service
where
  state = 'ENABLED';
```

# Table:  gcp_iam_role

An IAM role is an IAM entity that defines a set of permissions for making AWS service requests.

## Examples

### IAM role basic info

```sql
select
  name,
  role_id,
  deleted,
  description,
  title
from
  gcp_iam_role;
```


### List of IAM roles which are in BETA stage

```sql
select
  name,
  description,
  stage
from
  gcp_iam_role
where
  stage = 'BETA';
```
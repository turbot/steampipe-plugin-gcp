# Table: gcp_cloud_identity_group_membership

A Membership defines a relationship between a Group and an entity belonging to that Group, referred to as a "member".

**You must specify the identity group name** in the `where` clause (`where group_name=''`) to list the identity group memberships.

## Examples

### Basic info

```sql
select
  name,
  group_name,
  create_time,
  type,
  update_time
from
  gcp_cloud_identity_group_membership
where
  group_name = '123j0zll4288gmz';
```

### Get details of all google managed members in a group

```sql
select
  name,
  group_name,
  create_time,
  preferred_member_key ->> 'id' as member_id
from
  gcp_cloud_identity_group_membership
where
  group_name = '123j0zll4288gmz'
  and preferred_member_key ->> 'namespace' is null;
```

### Get all the groups that are members of a specific group

```sql
select
  name,
  group_name,
  create_time,
  preferred_member_key ->> 'id' as member_id
from
  gcp_cloud_identity_group_membership
where
  group_name = '123j0zll4288gmz'
  and type = 'GROUP';
```

### List roles assigned to each member of a group

```sql
select
  name,
  group_name,
  create_time,
  type,
  preferred_member_key ->> 'id' as member_id,
  role ->> 'name' as role_name,
  role -> 'expiryDetail' ->> 'expireTime' as role_expiry_time
from
  gcp_cloud_identity_group_membership,
  jsonb_array_elements(roles) as role
where
  group_name = '123j0zll4288gmz';
```

### Get details of a specific member of a group

```sql
select
  name,
  group_name,
  create_time,
  type,
  preferred_member_key ->> 'id' as member_id,
  role ->> 'name' as role_name,
  role -> 'expiryDetail' ->> 'expireTime' as role_expiry_time
from
  gcp_cloud_identity_group_membership,
  jsonb_array_elements(roles) as role
where
  group_name = '123j0zll4288gmz'
  and name = '123454620869324818189';
```
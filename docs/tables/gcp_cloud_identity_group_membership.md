---
title: "Steampipe Table: gcp_cloud_identity_group_membership - Query GCP Cloud Identity Group Memberships using SQL"
description: "Allows users to query GCP Cloud Identity Group Memberships, specifically the details of a group's members, providing insights into the group's structure and member roles."
folder: "Cloud Identity"
---

# Table: gcp_cloud_identity_group_membership - Query GCP Cloud Identity Group Memberships using SQL

Google Cloud Identity is a service within Google Cloud Platform that allows you to manage users, devices, and apps in a centralized manner. It provides a way to set up and manage memberships for various Google Cloud resources, including groups. Cloud Identity helps you stay informed about the structure and roles of your group memberships.

## Table Usage Guide

The `gcp_cloud_identity_group_membership` table provides insights into group memberships within Google Cloud Identity. As a system administrator, explore membership-specific details through this table, including member roles, member types, and associated metadata. Utilize it to uncover information about memberships, such as those with specific roles, the relationships between members, and the verification of member types.

**Important Notes**
- You must specify the identity group name in the `where` clause (`where group_name=''`) to list the identity group memberships.

## Examples

### Basic info
Explore the details of a specific group membership in Google Cloud Identity, focusing on its creation and update times. This query is useful in tracking changes and understanding the membership's history within a particular group.

```sql+postgres
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

```sql+sqlite
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
Explore which members of a specific group are managed by Google. This can be useful for understanding the management structure of your group and ensuring that all members are correctly managed.

```sql+postgres
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

```sql+sqlite
select
  name,
  group_name,
  create_time,
  json_extract(preferred_member_key, '$.id') as member_id
from
  gcp_cloud_identity_group_membership
where
  group_name = '123j0zll4288gmz'
  and json_extract(preferred_member_key, '$.namespace') is null;
```

### Get all the groups that are members of a specific group
Explore which groups are members of a specific group to understand their relationships and hierarchy. This is particularly useful for managing group memberships and assessing the structure within your GCP Cloud Identity.

```sql+postgres
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

```sql+sqlite
select
  name,
  group_name,
  create_time,
  json_extract(preferred_member_key, '$.id') as member_id
from
  gcp_cloud_identity_group_membership
where
  group_name = '123j0zll4288gmz'
  and type = 'GROUP';
```

### List roles assigned to each member of a group
Explore which roles are assigned to each member within a specific group, gaining insights into role distribution and expiry details. This helps in managing group permissions and understanding the access level of each member for security and administrative purposes.

```sql+postgres
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

```sql+sqlite
select
  g.name,
  g.group_name,
  g.create_time,
  g.type,
  json_extract(g.preferred_member_key, '$.id') as member_id,
  json_extract(role.value, '$.name') as role_name,
  json_extract(json_extract(role.value, '$.expiryDetail'), '$.expireTime') as role_expiry_time
from
  gcp_cloud_identity_group_membership as g,
  json_each(g.roles) as role
where
  g.group_name = '123j0zll4288gmz';
```

### Get details of a specific member of a group
This query is useful to gain insights into the specifics of a certain group member, such as their role and the expiry time of that role. It's particularly useful for managing roles and permissions within a group, ensuring the right access is provided to the right members at the right time.

```sql+postgres
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

```sql+sqlite
select
  g.name,
  g.group_name,
  g.create_time,
  g.type,
  json_extract(preferred_member_key, '$.id') as member_id,
  json_extract(role.value, '$.name') as role_name,
  json_extract(json_extract(role.value, '$.expiryDetail'), '$.expireTime') as role_expiry_time
from
  gcp_cloud_identity_group_membership,
  json_each(roles) as role
where
  group_name = '123j0zll4288gmz'
  and name = '123454620869324818189';
```
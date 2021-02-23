# Table: gcp_sql_database_instance

Cloud SQL instance is a database running in the cloud. It is used to store, replicate, and protect databases.

## Examples

### Basic info

```sql
select
  name,
  state,
  instance_type,
  database_version,
  machine_type,
  location
from
  gcp_sql_database_instance;
```

### List of users in the specified Cloud SQL instance.

```sql
select
  name,
  instance_users
from
  gcp_sql_database_instance
where
  name='my-sql-instance';
```

### List of replica databases and their master instances

```sql
select
  name,
  master_instance_name,
  replication_type,
  gce_zone as replica_database_zone
from
  gcp_sql_database_instance
where
  database_replication_enabled;
```

### List of assigned IP addresses to the database instances

```sql
select
  name,
  ip ->> 'ipAddress' as ip_address,
  ip ->> 'type' as type
from
  gcp_sql_database_instance,
  jsonb_array_elements(ip_addresses) as ip;
```

### List of external networks that can connect to the database instance

```sql
select
  name as instance_name,
  i ->> 'name' as authorized_network_name,
  i ->> 'value' as authorized_network_value,
  ip_configuration ->> 'ipv4Enabled' as ipv4_enabled
from
  gcp_sql_database_instance,
  jsonb_array_elements(ip_configuration -> 'authorizedNetworks') as i;
```

### List of database instances without application tag key

```sql
select
  name,
  tags
from
  gcp_sql_database_instance
where
  not tags :: JSONB ? 'application';
```

### Count of database instances per location

```sql
select
  location,
  count(*) instance_count
from
  gcp_sql_database_instance
group by
  location;
```

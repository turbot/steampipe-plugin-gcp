# Table: gcp_redis_instance

In Memorystore for Redis, an instance refers to an in-memory Redis data store. Memorystore for Redis Instances are also referred to as Redis instances. When you create a new Memorystore for Redis instance, you are creating a new Redis data store.

## Examples

### Basic info

```sql
select
  name,
  display_name,
  create_time,
  location_id,
  memory_size_gb,
  reserved_ip_range
from
  gcp_redis_instance;
```

### List instances that have authentication enabled

```sql
select
  name,
  display_name,
  create_time,
  location_id,
  memory_size_gb,
  reserved_ip_range
from
  gcp_redis_instance
where
  auth_enabled;
```

### List instances created in the last 7 days

```sql
select
  name,
  display_name,
  create_time,
  location_id,
  memory_size_gb,
  reserved_ip_range
from
  gcp_redis_instance
where
  create_time >= current_timestamp - interval '7 days';
```

### List the node details of each instance

```sql
select
  name,
  display_name,
  create_time,
  location_id,
  jsonb_pretty(nodes) as instance_nodes
from
  gcp_redis_instance
where
  name = 'instance-test'
  and location_id = 'europe-west3-c';
```

### List instances encrypted with customer-managed keys

```sql
select
  name,
  display_name,
  create_time,
  location_id,
  memory_size_gb,
  reserved_ip_range
from
  gcp_redis_instance
where
  customer_managed_key is not null;
```

### List instances that have transit mode disabled

```sql
select
  name,
  display_name,
  create_time,
  location_id,
  memory_size_gb,
  reserved_ip_range
from
  gcp_redis_instance
where
  transit_encryption_mode = 2;
```

### List the maintenance details of instances

```sql
select
  name,
  display_name,
  create_time,
  location_id,
  maintenance_policy,
  maintenance_schedule,
  maintenance_version,
  available_maintenance_versions
from
  gcp_redis_instance;
```

### List instances with direct peering access

```sql
select
  name,
  display_name,
  create_time,
  location_id,
  memory_size_gb,
  reserved_ip_range
from
  gcp_redis_instance
where
  connect_mode = 1;
```

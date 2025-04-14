---
title: "Steampipe Table: gcp_sql_database_instance - Query Google Cloud SQL Database Instances using SQL"
description: "Allows users to query Google Cloud SQL Database Instances, providing insights into the configuration, status, and performance of each instance."
folder: "SQL"
---

# Table: gcp_sql_database_instance - Query Google Cloud SQL Database Instances using SQL

Google Cloud SQL is a fully-managed database service that makes it easy to set up, maintain, manage, and administer your relational databases on Google Cloud Platform. It provides a cloud-based environment for running MySQL, PostgreSQL, and SQL Server databases. Google Cloud SQL offers high performance, scalability, and convenience for application developers.

## Table Usage Guide

The `gcp_sql_database_instance` table provides insights into the configuration and performance of Google Cloud SQL Database Instances. As a database administrator or developer, you can use this table to explore details about each instance, including its configuration, status, and performance metrics. This can help you optimize database performance, ensure proper configuration, and monitor the overall health of your databases.

## Examples

### Basic info
Explore which Google Cloud Platform SQL databases are currently active, their versions, and locations. This can help in understanding the distribution and usage of different databases across various regions.

```sql+postgres
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

```sql+sqlite
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
Assess the elements within a specific Cloud SQL instance to identify all associated users. This is beneficial in managing access control and maintaining security protocols.

```sql+postgres
select
  name,
  instance_users
from
  gcp_sql_database_instance
where
  name='my-sql-instance';
```

```sql+sqlite
select
  name,
  instance_users
from
  gcp_sql_database_instance
where
  name='my-sql-instance';
```

### List of replica databases and their master instances
Discover the segments that utilize replica databases by identifying their corresponding master instances. This can be beneficial in understanding the structure and distribution of your database system, particularly in scenarios where redundancy or load balancing is a key concern.

```sql+postgres
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

```sql+sqlite
select
  name,
  master_instance_name,
  replication_type,
  gce_zone as replica_database_zone
from
  gcp_sql_database_instance
where
  database_replication_enabled = 1;
```

### List of assigned IP addresses to the database instances
Explore which IP addresses have been assigned to your database instances. This can help you maintain a secure network and monitor potential unauthorized access.

```sql+postgres
select
  name,
  ip ->> 'ipAddress' as ip_address,
  ip ->> 'type' as type
from
  gcp_sql_database_instance,
  jsonb_array_elements(ip_addresses) as ip;
```

```sql+sqlite
select
  name,
  json_extract(ip.value, '$.ipAddress') as ip_address,
  json_extract(ip.value, '$.type') as type
from
  gcp_sql_database_instance,
  json_each(ip_addresses) as ip;
```

### List of external networks that can connect to the database instance
Explore which external networks have access to your database instance. This is useful to maintain security by ensuring only authorized networks can connect.

```sql+postgres
select
  name as instance_name,
  i ->> 'name' as authorized_network_name,
  i ->> 'value' as authorized_network_value,
  ip_configuration ->> 'ipv4Enabled' as ipv4_enabled
from
  gcp_sql_database_instance,
  jsonb_array_elements(ip_configuration -> 'authorizedNetworks') as i;
```

```sql+sqlite
select
  name as instance_name,
  json_extract(i.value, '$.name') as authorized_network_name,
  json_extract(i.value, '$.value') as authorized_network_value,
  json_extract(ip_configuration, '$.ipv4Enabled') as ipv4_enabled
from
  gcp_sql_database_instance,
  json_each(json_extract(ip_configuration, '$.authorizedNetworks')) as i;
```

### List of database instances without application tag key
Identify instances where database instances lack an application tag key. This is useful in understanding and rectifying configurations that are missing vital tagging, thereby improving resource management and organization.

```sql+postgres
select
  name,
  tags
from
  gcp_sql_database_instance
where
  not tags :: JSONB ? 'application';
```

```sql+sqlite
select
  name,
  tags
from
  gcp_sql_database_instance
where
  not json_valid(tags) or json_extract(tags, '$.application') is null;
```

### Count of database instances per location
Explore which locations have the highest number of database instances. This can help in understanding the distribution of resources and potentially identifying areas for infrastructure optimization.

```sql+postgres
select
  location,
  count(*) instance_count
from
  gcp_sql_database_instance
group by
  location;
```

```sql+sqlite
select
  location,
  count(*) instance_count
from
  gcp_sql_database_instance
group by
  location;
```
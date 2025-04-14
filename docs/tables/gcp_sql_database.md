---
title: "Steampipe Table: gcp_sql_database - Query Google Cloud SQL Databases using SQL"
description: "Allows users to query Google Cloud SQL Databases, providing detailed information about the database's configuration, status, and associated instances."
folder: "SQL"
---

# Table: gcp_sql_database - Query Google Cloud SQL Databases using SQL

Google Cloud SQL is a fully-managed database service that makes it easy to set up, maintain, manage, and administer relational databases on Google Cloud. It provides a cloud-based alternative to traditional on-premise databases, supporting both MySQL and PostgreSQL. Google Cloud SQL Databases offer high performance, scalability, and convenience.

## Table Usage Guide

The `gcp_sql_database` table provides insights into databases within Google Cloud SQL. As a database administrator or cloud engineer, you can explore database-specific details through this table, including the database name, instance, charset, and collation. Use it to manage and monitor your SQL databases efficiently, ensuring optimal performance and security.

## Examples

### Basic info
Explore the character set and collation configurations of your Google Cloud Platform SQL databases to ensure they align with your data encoding and sorting requirements. This can help maintain data integrity and prevent potential issues related to character representation and comparison.

```sql+postgres
select
  name,
  instance_name,
  charset,
  collation
from
  gcp_sql_database;
```

```sql+sqlite
select
  name,
  instance_name,
  charset,
  collation
from
  gcp_sql_database;
```

### Get the SQL Server version with which the database is to be made compatible
Determine the compatibility level of your database with different versions of SQL Server. This can be useful for planning version upgrades or ensuring backward compatibility with older versions.

```sql+postgres
select
  name,
  sql_server_database_compatibility_level
from
  gcp_sql_database;
```

```sql+sqlite
select
  name,
  sql_server_database_compatibility_level
from
  gcp_sql_database;
```

### Count of databases per instance
Analyze the settings to understand the distribution of databases across different instances. This can help in assessing the load distribution and managing resources more effectively.

```sql+postgres
select
  instance_name,
  count(*) as database_count
from
  gcp_sql_database
group by
  instance_name;
```

```sql+sqlite
select
  instance_name,
  count(*) as database_count
from
  gcp_sql_database
group by
  instance_name;
```
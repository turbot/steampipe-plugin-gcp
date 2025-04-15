---
title: "Steampipe Table: gcp_firestore_database - Query Google Cloud Firestore Databases using SQL"
description: "Allows users to query Google Cloud Firestore Databases, providing insights into the configuration and status of Firestore databases within a project."
---

# Table: gcp_firestore_database - Query Google Cloud Firestore Databases using SQL

Firestore is a NoSQL document database built for automatic scaling, high performance, and ease of application development. This table provides information about Firestore databases in your GCP project.

## Table Usage Guide

The `gcp_firestore_database` table provides insights into Firestore Databases within Google Cloud. As a database administrator or cloud engineer, you can explore database-specific details through this table, including the database name, type, and settings. Use it to manage and monitor your SQL databases efficiently, ensuring optimal performance and security.

## Examples

### Basic info
Explore the type and concurrency mode of your Google Cloud Platform Firestore databases to ensure they align with your requirements. This can help prevent potential issues related to contention when querying data.

```sql+postgres
select
  name,
  uid,
  type,
  location,
  concurrency_mode,
  create_time
from
  gcp_firestore_database;```

```sql+sqlite
select
  name,
  uid,
  type,
  location,
  concurrency_mode,
  create_time
from
  gcp_firestore_database;
```

### Get details of a specific database
Determine the versioning settings and status of your database. This can be useful for ensuring incorrect updates can be rolled back easily.

```sql+postgres
select
  name,
  uid,
  type,
  location,
  create_time,
  version_retention_period,
  earliest_version_time
from
  gcp_firestore_database
where
  title = '(default)';
```

```sql+sqlite
select
  name,
  uid,
  type,
  location,
  create_time,
  version_retention_period,
  earliest_version_time
from
  gcp_firestore_database
where
  title = '(default)';
```

### List databases by type
Determine the versioning settings and status of your database. This can be useful for ensuring incorrect updates can be rolled back easily.

```sql+postgres
select
  name,
  uid,
  type,
  location
from
  gcp_firestore_database
where
  type = 'FIRESTORE_NATIVE';
```

```sql+sqlite
select
  name,
  uid,
  type,
  location
from
  gcp_firestore_database
where
  type = 'FIRESTORE_NATIVE';
```

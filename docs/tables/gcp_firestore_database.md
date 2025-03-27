# Table: gcp_firestore_database

Firestore is a NoSQL document database built for automatic scaling, high performance, and ease of application development. This table provides information about Firestore databases in your GCP project.

## Examples

### List all Firestore databases
```sql
select
  name,
  uid,
  type,
  location,
  create_time
from
  gcp_firestore_database;
```

### Get details of a specific database
```sql
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
```sql
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

### List databases by location
```sql
select
  name,
  uid,
  type,
  location
from
  gcp_firestore_database
where
  location = 'eur3';

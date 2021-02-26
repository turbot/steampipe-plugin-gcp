# Table: gcp_sql_database

Cloud SQL is a fully-managed database service that helps to set up, maintain, manage, and administer relational databases on Google Cloud Platform.

## Examples

### Basic info

```sql
select
  name,
  instance_name,
  charset,
  collation
from
  gcp_sql_database;
```


### Get the SQL Server version with which the database is to be made compatible

```sql
select
  name,
  sql_server_database_compatibility_level
from
  gcp_sql_database;
```


### Count of databases per instance

```sql
select
  instance_name,
  count(*) as database_count
from
  gcp_sql_database
group by
  instance_name;
```
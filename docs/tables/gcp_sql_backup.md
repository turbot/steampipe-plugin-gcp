---
title: "Steampipe Table: gcp_sql_backup - Query Google Cloud SQL Backups using SQL"
description: "Allows users to query Google Cloud SQL Backups, specifically the backup details, providing insights into backup configurations, timings, and statuses."
folder: "SQL"
---

# Table: gcp_sql_backup - Query Google Cloud SQL Backups using SQL

Google Cloud SQL Backups is a service within Google Cloud that allows you to create, configure, and manage backups for your SQL databases. It provides a streamlined way to safeguard your data and maintain business continuity. Google Cloud SQL Backups helps you ensure data integrity, recover from disasters, and meet compliance requirements.

## Table Usage Guide

The `gcp_sql_backup` table provides insights into backups within Google Cloud SQL. As a database administrator, explore backup-specific details through this table, including backup configurations, timings, and statuses. Utilize it to uncover information about backups, such as their configurations, the timing of the last successful backup, and the status of ongoing backups.

## Examples

### Basic info
Determine the status and timing details of your Google Cloud Platform SQL backups. This can help you understand backup health and scheduling, essential for maintaining data integrity and planning recovery scenarios.

```sql+postgres
select
  id,
  instance_name,
  description,
  status,
  end_time,
  enqueued_time,
  start_time,
  window_start_time
from
  gcp_sql_backup;
```

```sql+sqlite
select
  id,
  instance_name,
  description,
  status,
  end_time,
  enqueued_time,
  start_time,
  window_start_time
from
  gcp_sql_backup;
```

### Count of backups by their type (i.e AUTOMATED and ON_DEMAND)
Determine the distribution of backup types to understand your database's backup strategy. This can help in assessing the balance between automated and on-demand backups, and optimize your data protection approach.

```sql+postgres
select
  type,
  count(*) as backup_count
from
  gcp_sql_backup
group by
  type;
```

```sql+sqlite
select
  type,
  count(*) as backup_count
from
  gcp_sql_backup
group by
  type;
```

### Get the error message if the backup failed
Determine the areas in which a backup failure has occurred in your Google Cloud Platform SQL database. This query can be used to identify the specific instances and error messages associated with each failure, helping you troubleshoot and resolve issues more efficiently.

```sql+postgres
select
  id,
  instance_name,
  e ->> 'code' as error_code,
  e ->> 'message' as error_message
from
  gcp_sql_backup,
  jsonb_array_elements(error) as e
where
  status = 'FAILED';
```

```sql+sqlite
select
  b.id,
  b.instance_name,
  json_extract(e.value, '$.code') as error_code,
  json_extract(e.value, '$.message') as error_message
from
  gcp_sql_backup as b,
  json_each(error) as e
where
  status = 'FAILED';
```
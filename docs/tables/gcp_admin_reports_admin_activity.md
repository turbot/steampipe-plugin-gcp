---
title: "Steampipe Table: gcp_admin_reports_admin_activity - Query GCP Admin Reports Admin Activity Events using SQL"
description: "Allows users to query admin activity events from the GCP Admin Reports API, providing insights into administrative actions and login events."
folder: "Cloud Admin Reports"
---

# Table: gcp_admin_reports_admin_activity - Query GCP Admin Reports Admin Activity Events using SQL

Google Admin Reports Admin Activity provides audit logs for administrative actions and login events performed in your GCP/Google Workspace domain. These records enable you to detect privilege or user modifications, to investigate security events, and to monitor compliance.

## Table Usage Guide

The `gcp_admin_reports_admin_activity` table is ideal for tracking admin operations such as login attempts, password changes, and role assignments. Use it to monitor and audit critical administrative events.

> :point_right: Notice that the event_name are inside brackets, it's because we can have several events for the same entry, example : `[CHANGE_PASSWORD CHANGE_PASSWORD_ON_NEXT_LOGIN]`

## Examples

### 1. List recent admin activities

Retrieve admin events in the last 24 hours:

```sql
select
  time,
  actor_email,
  event_name,
  ip_address
from
  gcp_admin_reports_admin_activity
where
  time > now() - '24 hours';
```

### 2. Filter by event name

Show all changes of password:

```sql
select
  time,
  actor_email,
  event_name,
  user_email,
  ip_address
from 
  gcp_admin_reports_admin_activity
where 
  event_name like '%CHANGE_PASSWORD%';
```

### 3. Identify failed login attempts

Find admin login events originating from an unexpected IP range:

```sql
select
  time,
  actor_email,
  ip_address,
  event_name
from
  gcp_admin_reports_admin_activity
where
  event_name = 'login_failure'
  and ip_address >= '203.0.113.0'
  and ip_address <= '203.0.113.255';
```

### 4. Get activities within a custom time window

Query admin activities between two timestamps:

```sql
select
  time,
  actor_email,
  event_name,
  unique_qualifier
from
  gcp_admin_reports_admin_activity
where
  time between '2025-06-15T00:00:00Z' and '2025-06-20T23:59:59Z';
```

### 5. Top admin users by activity count

Aggregate total admin events per user in the last month:

```sql
select
  actor_email,
  count(*) as total_events
from
  gcp_admin_reports_admin_activity
where
  time > now() - '1 month'::interval
group by
  actor_email
order by
  total_events desc
limit 10;
```
---
title: "Steampipe Table: gcp_admin_reports_admin_activity - Query GCP Admin Reports Admin Activity Events using SQL"
description: "Allows users to query admin activity events from the GCP Admin Reports API, providing insights into administrative actions and login events."
folder: "Cloud Admin Reports"
---

# Table: gcp_admin_reports_admin_activity - Query GCP Admin Reports Admin Activity Events using SQL

Google Admin Reports Admin Activity provides audit logs for administrative actions and login events performed in your GCP/Google Workspace domain. These records enable you to detect privilege or user modifications, to investigate security events, and to monitor compliance.

## Table Usage Guide

The `gcp_admin_reports_admin_activity` table is ideal for tracking admin operations such as user creation, password changes, and role assignments. Use it to monitor and audit critical administrative events.

> :point_right: Notice that the event_name are inside brackets, it's because we can have several events for the same entry, example : `[CHANGE_PASSWORD CHANGE_PASSWORD_ON_NEXT_LOGIN]`

## Examples

### Basic info

Retrieve events performed by administrators of your GCP/Google Workspace domain in the last 24 hours.

```sql
select
  time,
  actor_email,
  event_name,
  ip_address,
  events
from
  gcp_admin_reports_admin_activity
where
  time > now() - '1 day'::interval;
```

### List all password change events

Show all changes of password performed by administrators on users.

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

### List all user creation events

Find user creation events performed by admins originating from an unexpected IP range.

```sql
select
  time,
  actor_email,
  ip_address,
  user_email,
  event_name
from
  gcp_admin_reports_admin_activity
where
  event_name = '[CREATE_USER]'
  and ip_address >= '203.0.113.0'
  and ip_address <= '203.0.113.255';
```

### Get activities within a custom time window

Query admin activities between two timestamps.

```sql
select
  time,
  actor_email,
  event_name,
  events
from
  gcp_admin_reports_admin_activity
where
  time between '2025-06-15T00:00:00Z' and '2025-06-20T23:59:59Z';
```
---
title: "Steampipe Table: gcp_admin_reports_mobile_activity - Query GCP Admin Reports Mobile Activity Events using SQL"
description: "Allows users to query mobile activity events from the GCP Admin Reports API, providing insights into device usage and mobile access patterns."
folder: "Cloud Admin Reports"
---

# Table: gcp_admin_reports_mobile_activity - Query GCP Admin Reports Mobile Activity Events using SQL

Google Admin Reports Mobile Activity captures events related to mobile device actions—such as mobile logins, syncs, and policy enforcements—performed by users in your organization. This table surfaces those audit records via Steampipe’s SQL interface.

## Table Usage Guide

Use the `gcp_admin_reports_mobile_activity` table to investigate how users interact with Google services from their mobile devices. Track device enrollments, removals, sync operations, and any policy-related events.

## Examples

### 1. Recent mobile logins

Retrieve mobile login events in the last 6 hours:

```sql
select
  time,
  actor_email,
  event_name,
  device_model
from
  gcp_admin_reports_mobile_activity
where
  time > now() - '6 hours'::interval;
```

### 2. Device sync events for a user

Show sync operations by [alice@example.com] over the past 3 days:

```sql
select
  time,
  actor_email,
  event_name,
  device_model
from
  gcp_admin_reports_mobile_activity
where
  actor_email = 'alice@example.com'
  and event_name like '%DEVICE_SYNC_EVENT%'
  and time > now() - '3 days'::interval;
```

### 3. Connections from a new device

Identify all connections from a new device in the last week:

```sql
select
  time,
  actor_email,
  event_name, device_id,
  device_model
from
  gcp_admin_reports_mobile_activity
where
  event_name = '[DEVICE_REGISTER_UNREGISTER_EVENT]';
```

### 4. Custom time window analysis

Query mobile activities between two timestamps:

```sql
select
  time,
  actor_email,
  event_name, device_id,
  device_model
from
  gcp_admin_reports_mobile_activity
where
  time between '2025-06-10T00:00:00Z' and '2025-06-15T23:59:59Z';
```

### 5. Top users by mobile event count

Aggregate total mobile events per user in the last 24 hours:

```sql
select
  actor_email,
  count(*) as total_events
from
  gcp_admin_reports_mobile_activity
where
  time > now() - '24 hours'::interval
group by
  actor_email
order by
  total_events desc
limit 10;
```

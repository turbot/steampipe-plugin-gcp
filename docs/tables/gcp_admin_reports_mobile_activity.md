---
title: "Steampipe Table: gcp_admin_reports_mobile_activity - Query GCP Admin Reports Mobile Activity Events using SQL"
description: "Allows users to query mobile activity events from the GCP Admin Reports API, providing insights into device usage and mobile access patterns."
folder: "Cloud Admin Reports"
---

# Table: gcp_admin_reports_mobile_activity - Query GCP Admin Reports Mobile Activity Events using SQL

Google Admin Reports Mobile Activity captures events related to device actions—such as logins, syncs, updates of the OS and usage of a new device by members of your organization.

## Table Usage Guide

Use the `gcp_admin_reports_mobile_activity` table to investigate how users interact with Google services from their devices. Track device enrollments, removals, sync operations, and any policy-related events.

## Examples

### Basic info

Retrieve device-related events in the last 24 hours.

```sql
select
  time,
  actor_email,
  event_name,
  device_model
from
  gcp_admin_reports_mobile_activity
where
  time > now() - '1 day'::interval;
```

### Device sync events for a user

Show sync operations by [alice@example.com] over the past 3 days.

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

### Connections from a new device

Identify all connections from a new device in the last week.

```sql
select
  time,
  actor_email,
  event_name, device_id,
  device_model
from
  gcp_admin_reports_mobile_activity
where
  event_name = '[DEVICE_REGISTER_UNREGISTER_EVENT]'
  and time > now() - '1 week'::interval;
```

### Custom time window analysis

Query device activities between two timestamps.

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
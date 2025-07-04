---
title: "Steampipe Table: gcp_admin_reports_token_activity - Query GCP Admin Reports Token Activity Events using SQL"
description: "Allows users to query token activity events from the GCP Admin Reports API, providing insights into OAuth and API token usage and revocation events."
folder: "Cloud Admin Reports"
---

# Table: gcp_admin_reports_token_activity - Query GCP Admin Reports Token Activity Events using SQL

Google Admin Reports Token Activity captures events related to OAuth and API tokens—such as token authorization (connections to other services using Google account) and revocation. 

## Table Usage Guide

Use the `gcp_admin_reports_token_activity` table to monitor token lifecycle events, detect unauthorized token usage, and audit application integrations.

## Examples

### Basic info

Retrieve token activity events in the last 24 hours.

```sql
select
  time,
  actor_email,
  event_name,
  app_name
from
  gcp_admin_reports_token_activity
where
  time > now() - '1 day'::interval;
```

### Token activity related to a specific app

Identify event related to the Google Chrome app over the last week.

```sql
select
  time,
  actor_email,
  event_name
from
  gcp_admin_reports_token_activity
where
  app_name = 'Google Chrome'
  and time > now() - '7 days'::interval;
```

### Token revocations for a user

Show all token deletion events by [alice@example.com] in the past 2 days.

```sql
select
  time,
  event_name,
  app_name
from
  gcp_admin_reports_token_activity
where
  actor_email = 'alice@example.com'
  and event_name = '[revoke]'
  and time > now() - '2 days'::interval;
```

### Custom time window audit

Query token activity between two specific timestamps.

```sql
select
  time,
  actor_email,
  event_name,
  app_name
from
  gcp_admin_reports_token_activity
where
  time between '2025-06-01T00:00:00Z' and '2025-06-07T23:59:59Z';
```
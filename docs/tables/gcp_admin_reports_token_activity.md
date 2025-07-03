---
title: "Steampipe Table: gcp_admin_reports_token_activity - Query GCP Admin Reports Token Activity Events using SQL"
description: "Allows users to query token activity events from the GCP Admin Reports API, providing insights into OAuth and API token usage and revocation events."
folder: "Cloud Admin Reports"
---

# Table: gcp_admin_reports_token_activity - Query GCP Admin Reports Token Activity Events using SQL

Google Admin Reports Token Activity captures events related to OAuth and API tokens—such as token creation, refresh, revocation, and invalidation—within your organization. This table exposes those audit records via Steampipe’s SQL interface.

## Table Usage Guide

Use the `gcp_admin_reports_token_activity` table to monitor token lifecycle events, detect unauthorized token usage, and audit application integrations.

## Examples

### 1. Recent token creations

Retrieve token creation events in the last 8 hours:

```sql
select
  time,
  actor_email,
  event_name,
  app_name
from
  gcp_admin_reports_token_activity
where
  time > now() - '8 hours'::interval;
```

### 2. Token revocations for a user

Show all token deletion events by [alice@example.com] in the past 2 days:

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

### 3. Refresh token events

Identify event related to the Google Chrome app over the last week:

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

### 4. Custom time window audit

Query token activity between two specific timestamps:

```sql
select
  time,
  actor_email,
  event_name
from
  gcp_admin_reports_token_activity
where
  time between '2025-06-01T00:00:00Z' and '2025-06-07T23:59:59Z';
```

### 5. Top token event types in the last month

Aggregate counts of each token event type in the last 30 days:

```sql
select
  event_name,
  count(*) as event_count
from
  gcp_admin_reports_token_activity
where
  time > now() - '30 days'::interval
group by
  event_name
order by
  event_count desc;
```
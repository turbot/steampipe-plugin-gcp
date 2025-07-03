---
title: "Steampipe Table: gcp\_admin\_reports\_login\_activity - Query GCP Admin Reports Login Activity Events using SQL"
description: "Allows users to query login activity events from the GCP Admin Reports API, providing insights into user login behavior and authentication events."
folder: "Cloud Admin Reports"
---

# Table: gcp_admin_reports_login_activity - Query GCP Admin Reports Login Activity Events using SQL

Google Admin Reports Login Activity captures authentication events such as successful and failed logins, MFA completions, and session terminations within your GSuite domain. This table surfaces those audit records in an SQL interface via Steampipe.

## Table Usage Guide

The `gcp_admin_reports_login_activity` table is designed for monitoring user authentication actions. Use it to track login success/failure rates, identify anomalous access patterns, and support security investigations.

## Examples

### 1. List recent login events

Retrieve login events in the last 12 hours:

```sql
select
  time,
  actor_email,
  event_name,
  ip_address
from
  gcp_admin_reports_login_activity
where
  time > now() - '12 hours'::interval;
```

### 2. Filter by specific user failure events

Show all failed login attempts by [bob@example.com] in the last 7 days:

```sql
select
  time,
  event_name,
  ip_address
from
  gcp_admin_reports_login_activity
where
  actor_email = 'bob@example.com'
  and event_name = 'login_failure'
  and time > now() - '7 days'::interval;
```

### 3. Identify passwords changes

Find password change events across all users in the last week:

```sql
select
  time,
  actor_email,
  event_name
from
  gcp_admin_reports_login_activity
where
  event_name = 'password_edit'
  and time > now() - '7 days'::interval;
```

### 4. Custom time window analysis

Query login activities between two timestamps:

```sql
select
  time,
  actor_email,
  event_name
from
  gcp_admin_reports_login_activity
where
  time between '2025-06-10T00:00:00Z' and '2025-06-15T23:59:59Z';
```

### 5. Top IP addresses by login count

Identify the top source IPs initiating login events in the last month:

```sql
select
  ip_address,
  count(*) as login_count
from
  gcp_admin_reports_login_activity
where
  time >now() - '30 days'::interval
group by
  ip_address
order by
  login_count desc
limit 10;
```
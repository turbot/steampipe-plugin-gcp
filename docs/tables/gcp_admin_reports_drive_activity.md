---
title: "Steampipe Table: gcp_admin_reports_drive_activity - Query GCP Admin Reports Drive Activity Events using SQL"
description: "Allows users to query Drive activity events from the GCP Admin Reports API, providing insights into file operations and user interactions on Google Drive."
folder: "Cloud Admin Reports"
---

# Table: gcp_admin_reports_drive_activity - Query GCP Admin Reports Drive Activity Events using SQL

Google Admin Reports Drive Activity captures detailed events related to Google Drive operations—such as file views, edits, downloads, creations, and deletions—performed by users within your organization.

## Table Usage Guide

The `gcp_admin_reports_drive_activity` table lets administrators and security teams investigate Drive-based activities. Use it to track who accessed or modified files, when operations occurred, and contextual metadata such as IP address and event types.

> :point_right: Notice that the event_name are inside brackets, it's because we can have several events for the same entry, example : `[edit change_user_access add_to_folder upload]`
>
> :exclamation: For improved performance, it is advised that you use the option `time` to limit the result set to a specific time period.

## Examples

### Basic info
Retrieve events occuring in the Google Drive of your organization in the last 24 hours, showing user and file names.

```sql
select
  time,
  actor_email,
  file_name,
  event_name
from
  gcp_admin_reports_drive_activity
where
  time > now() - interval '1 day';
```

### Show events related to a specific file
Show Drive edits and views performed on the file Passwords.txt in the last week.

```sql
select
  time,
  actor_email,
  event_name,
  ip_address
from
  gcp_admin_reports_drive_activity
where
  file_name = 'Passwords.txt'
  and event_name in ('[edit]', '[view]')
  and time > now() - '1 week'::interval;
```

### Find activities from a specific IP address
Identify all Drive operations originating from a specific IP address.

```sql
select
  time,
  actor_email,
  event_name,
  file_name
from
  gcp_admin_reports_drive_activity
where
  ip_address = '8.8.8.8';
```

### Get events within a custom time window
Query Drive activities between two timestamps.

```sql
select
  time,
  actor_email,
  event_name,
  file_name
from
  gcp_admin_reports_drive_activity
where
  time between '2025-06-15T00:00:00Z' and '2025-06-16T23:59:59Z';
```

### Top users by activity count
Aggregate total Drive events per user in the last 5 hours.

```sql
select
  actor_email,
  count(*) as total_events
from
  gcp_admin_reports_drive_activity
where
  time >= now() - interval '5 hours'
group by
  actor_email
order by
  total_events desc
limit 10;
```
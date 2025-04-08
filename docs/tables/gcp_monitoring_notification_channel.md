---
title: "Steampipe Table: gcp_monitoring_notification_channel - Query GCP Monitoring Notification Channels using SQL"
description: "Allows users to query GCP Monitoring Notification Channels, specifically to retrieve channel-specific details such as type, display name, labels, and verification status."
folder: "Cloud Monitoring"
---

# Table: gcp_monitoring_notification_channel - Query GCP Monitoring Notification Channels using SQL

GCP Monitoring Notification Channels are a part of Google Cloud Monitoring that allows you to configure where and how you want to receive notifications when incidents occur. These channels can be configured to send notifications via email, SMS, or to third-party services. This enables you to stay informed about the health and performance of your GCP resources and take appropriate actions when predefined conditions are met.

## Table Usage Guide

The `gcp_monitoring_notification_channel` table provides insights into Monitoring Notification Channels within Google Cloud Platform. As a DevOps engineer or system administrator, explore channel-specific details through this table, such as the type of channel, display name, labels, and verification status. Utilize it to manage and monitor your notification channels effectively, ensuring you receive timely alerts about incidents in your GCP resources.

## Examples

### List of monitoring notification channel which are not verified
Determine the areas in which certain monitoring notification channels remain unverified. This query is useful for identifying potential gaps in your monitoring system, allowing for prompt verification and ensuring comprehensive coverage.

```sql+postgres
select
  name,
  display_name,
  type,
  verification_status
from
  gcp_monitoring_notification_channel
where
  verification_status <> 'VERIFIED';
```

```sql+sqlite
select
  name,
  display_name,
  type,
  verification_status
from
  gcp_monitoring_notification_channel
where
  verification_status <> 'VERIFIED' OR verification_status is null;
```

### List of monitoring notification channel which are not enabled
Explore which monitoring notification channels in your Google Cloud Platform are not currently enabled. This can help you identify potential gaps in your monitoring strategy and ensure that all necessary channels are active.

```sql+postgres
select
  name,
  display_name,
  enabled
from
  gcp_monitoring_notification_channel
where
  not enabled;
```

```sql+sqlite
select
  name,
  display_name,
  enabled
from
  gcp_monitoring_notification_channel
where
  enabled = 0;
```
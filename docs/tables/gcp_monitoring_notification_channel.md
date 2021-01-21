# Table:  gcp_monitoring_notification_channel

Monitoring notification channel manages the sending of notifications to Pub/Sub-based notification channels in this project.

## Examples

### List of monitoring notification channel which are not verified

```sql
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


### List of monitoring notification channel which are not enabled

```sql
select
  name,
  display_name,
  enabled
from
  gcp_monitoring_notification_channel
where
  not enabled;
```
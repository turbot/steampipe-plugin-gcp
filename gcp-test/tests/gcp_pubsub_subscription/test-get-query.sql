select name, topic, topic_name, ack_deadline_seconds, message_retention_duration, push_config_endpoint, project, location, labels
from gcp.gcp_pubsub_subscription
where name = '{{ resourceName }}'
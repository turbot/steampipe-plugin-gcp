select name, topic, ack_deadline_seconds, message_retention_duration, push_config_endpoint, project
from gcp.gcp_pubsub_subscription
where name = '{{resourceName}}'
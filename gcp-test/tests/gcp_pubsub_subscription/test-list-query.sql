select name, topic, ack_deadline_seconds
from gcp.gcp_pubsub_subscription
where name = '{{resourceName}}'
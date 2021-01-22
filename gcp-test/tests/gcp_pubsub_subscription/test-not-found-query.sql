select name, topic
from gcp.gcp_pubsub_subscription
where name = 'dummy-{{resourceName}}'
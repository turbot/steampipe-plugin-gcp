select title, tags, akas
from gcp.gcp_pubsub_subscription
where name = '{{resourceName}}'
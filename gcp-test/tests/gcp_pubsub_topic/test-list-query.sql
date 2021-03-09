select name, tags
from gcp.gcp_pubsub_topic
where name = '{{resourceName}}'
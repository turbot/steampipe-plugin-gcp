select title, tags, akas
from gcp.gcp_pubsub_topic
where name = '{{resourceName}}'
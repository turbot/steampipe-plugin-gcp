select name, project
from gcp.gcp_pubsub_topic
where name = '{{resourceName}}'
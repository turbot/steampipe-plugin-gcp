select name, project, location, labels
from gcp.gcp_pubsub_topic
where name = '{{ resourceName }}'
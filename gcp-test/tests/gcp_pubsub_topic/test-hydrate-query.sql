select name, iam_policy, project
from gcp.gcp_pubsub_topic
where name = '{{resourceName}}'
select name, iam_policy, project
from gcp.gcp_pubsub_subscription
where name = '{{resourceName}}'
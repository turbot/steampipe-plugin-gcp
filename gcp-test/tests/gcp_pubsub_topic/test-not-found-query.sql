select name, kms_key_name
from gcp.gcp_pubsub_topic
where name = 'dummy-{{resourceName}}'
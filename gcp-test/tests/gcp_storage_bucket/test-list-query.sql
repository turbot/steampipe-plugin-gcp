select name, iam_policy, acl, default_object_acl, location_type, storage_class, billing_requester_pays, default_kms_key_name, log_bucket, versioning_enabled, website_main_page_suffix, website_not_found_page, cors, lifecycle_rules, retention_policy
from gcp.gcp_storage_bucket
where akas::text = '["{{ output.resource_aka.value }}"]'
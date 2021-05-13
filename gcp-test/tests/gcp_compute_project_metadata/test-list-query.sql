select name, default_service_account, common_instance_metadata->'items' as metadata_items, common_instance_metadata->'kind' as metadata_kind, default_network_tier, enabled_features, kind, self_link
from gcp.gcp_compute_project_metadata
where akas::text = '["{{ output.project_aka.value }}"]';
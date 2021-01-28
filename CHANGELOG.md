# CHANGELOG for GCP Plugin for Steampipe

## v0.0.4 [2021-01-28]

_What's new?_

- Added: `gcp_cloudfunctions_function` table
- Added: `gcp_compute_global_address` table
- Added: `gcp_compute_global_forwarding_rule` table
- Added: `gcp_compute_instance` table
- Added: `gcp_storage_bucket` table

- Updated: `gcp_iam_role` table. Added `is_gcp_managed` boolean field to distinguish between GCP managed and Customer managed roles.

_Bug fixes_

- Fixed: `gcp_iam_role` table. Updated `included_permissions` field to have details of role grants for list call.

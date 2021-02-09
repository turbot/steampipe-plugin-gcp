## v0.0.5 [2021-02-24]

_What's new?_

- New tables added to plugin

  - gcp_compute_address ([#29](https://github.com/turbot/steampipe-plugin-gcp/issues/29))
  - gcp_compute_disk ([#47](https://github.com/turbot/steampipe-plugin-gcp/issues/47))
  - gcp_compute_firewall ([#42](https://github.com/turbot/steampipe-plugin-gcp/issues/42))
  - gcp_compute_forwarding_rule ([#53](https://github.com/turbot/steampipe-plugin-gcp/issues/53))
  - gcp_compute_network ([#43](https://github.com/turbot/steampipe-plugin-gcp/issues/43))
  - gcp_compute_router ([#51](https://github.com/turbot/steampipe-plugin-gcp/issues/51))
  - gcp_compute_snapshot ([#60](https://github.com/turbot/steampipe-plugin-gcp/issues/60))

_Enhancements_

- Added field `location` to resource tables that are not regional with value as `global`

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



## v0.2.0 [2021-02-25]

_What's new?_

- New tables added
  - gcp_sql_backup
  - gcp_sql_database
  - gcp_sql_database_instance

_Bug fixes_

- Updated `gcp_compute_instance` table `network_tags` field to display value correctly ([#114](https://github.com/turbot/steampipe-plugin-gcp/pull/114))

## v0.1.1 [2021-02-22]

_Bug fixes_

- Now union query for multiple projects will work if different credential files are used for project connection in config ([#116](https://github.com/turbot/steampipe-plugin-gcp/issues/116))

- Updated `gcp_storage_bucket` table `labels` field to display value correctly ([#115](https://github.com/turbot/steampipe-plugin-gcp/issues/115))

## v0.1.0 [2021-02-18]

_What's new?_

- Added support for [connection configuration](https://github.com/turbot/steampipe-plugin-gcp/blob/main/docs/index.md#connection-configuration). You may specify gcp `project` and `credential_file` for each connection in a configuration file. You can have multiple gcp connections, each configured for a different gcp project.

- If the project id to connect to is not specified in connection configuration file or through `CLOUDSDK_CORE_PROJECT` environment variable. Now plugin will use active project, as returned by the `gcloud config get-value project` command.

_Enhancements_

- Added `location` column to `gcp_compute_image`, `gcp_compute_snapshot` and `gcp_monitoring_notification_channel`, `gcp_pubsub_snapshot`, `gcp_pubsub_subscription` and `gcp_pubsub_topic` tables.
- Added `iamPolicy` column to `gcp_compute_instance` table.
- Added `disabled` and `oauth2_client_id` columns to gcp_service_account table.

## v0.0.6 [2021-02-11]

_What's new?_

- New tables added to plugin

  - [gcp_compute_backend_bucket](https://github.com/turbot/steampipe-plugin-gcp/blob/main/docs/tables/gcp_compute_backend_bucket.md) ([#70](https://github.com/turbot/steampipe-plugin-gcp/issues/70))
  - [gcp_compute_backend_service](https://github.com/turbot/steampipe-plugin-gcp/blob/main/docs/tables/gcp_compute_backend_service.md) ([#80](https://github.com/turbot/steampipe-plugin-gcp/issues/80))
  - [gcp_compute_image](https://github.com/turbot/steampipe-plugin-gcp/blob/main/docs/tables/gcp_compute_image.md) ([#45](https://github.com/turbot/steampipe-plugin-gcp/issues/45))
  - [gcp_compute_instance_template](https://github.com/turbot/steampipe-plugin-gcp/blob/main/docs/tables/gcp_compute_instance_template.md) ([#84](https://github.com/turbot/steampipe-plugin-gcp/issues/84))
  - [gcp_compute_node_group](https://github.com/turbot/steampipe-plugin-gcp/blob/main/docs/tables/gcp_compute_node_group.md) ([#58](https://github.com/turbot/steampipe-plugin-gcp/issues/58))
  - [gcp_compute_node_template](https://github.com/turbot/steampipe-plugin-gcp/blob/main/docs/tables/gcp_compute_node_template.md) ([#87](https://github.com/turbot/steampipe-plugin-gcp/issues/87))
  - [gcp_compute_subnetwork](https://github.com/turbot/steampipe-plugin-gcp/blob/main/docs/tables/gcp_compute_subnetwork.md) ([#68](https://github.com/turbot/steampipe-plugin-gcp/issues/68))
  - [gcp_compute_target_pool](https://github.com/turbot/steampipe-plugin-gcp/blob/main/docs/tables/gcp_compute_target_pool.md) ([#81](https://github.com/turbot/steampipe-plugin-gcp/issues/81))
  - [gcp_compute_target_vpn_gateway](https://github.com/turbot/steampipe-plugin-gcp/blob/main/docs/tables/gcp_compute_target_vpn_gateway.md) ([#65](https://github.com/turbot/steampipe-plugin-gcp/issues/65))
  - [gcp_compute_url_map](https://github.com/turbot/steampipe-plugin-gcp/blob/main/docs/tables/gcp_compute_url_map.md) ([#85](https://github.com/turbot/steampipe-plugin-gcp/issues/85))
  - [gcp_compute_vpn_tunnel](https://github.com/turbot/steampipe-plugin-gcp/blob/main/docs/tables/gcp_compute_vpn_tunnel.md) ([#63](https://github.com/turbot/steampipe-plugin-gcp/issues/63))

## v0.0.5 [2021-02-04]

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

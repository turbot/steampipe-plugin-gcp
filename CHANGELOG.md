## v1.6.0 [2025-04-14]

_Enhancements_

- Added `folder` metadata to the documentation of all the GCP tables for improved organization on the Steampipe Hub. ([#740](https://github.com/turbot/steampipe-plugin-gcp/pull/740))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.11.5](https://github.com/turbot/steampipe-plugin-sdk/releases/tag/v5.11.5) ([#737](https://github.com/turbot/steampipe-plugin-gcp/pull/737))

## v1.5.0 [2025-02-15]

_Dependencies_

- Recompiled plugin with `Go v1.231.1`. ([#714](https://github.com/turbot/steampipe-plugin-gcp/pull/714))
- Recompiled plugin with [steampipe-plugin-sdk v5.11.3](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5113-2025-02-11) ([#714](https://github.com/turbot/steampipe-plugin-gcp/pull/714))

## v1.4.0 [2025-02-12]

_Enhancements_

- Added column `billing_information` to `gcp_organization_project` and `gcp_project` tables. ([#665](https://github.com/turbot/steampipe-plugin-gcp/pull/665)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)

_Bug fixes_

- Fixed the `gcp_cloud_run_job` table to correctly return data from the service supported regions instead of returning a `permission denied` error. ([#713](https://github.com/turbot/steampipe-plugin-gcp/pull/713))

_Deprecated_

- The `project` column in the `gcp_billing_account` table is now deprecated and will be removed in the next major version. ([#665](https://github.com/turbot/steampipe-plugin-gcp/pull/665)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)

## v1.3.0 [2025-02-05]

_What's new?_

- New tables added
 - [gcp_redis_cluster](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_redis_cluster) ([#679](https://github.com/turbot/steampipe-plugin-gcp/pull/679)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)
 - [gcp_organization_project](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_organization_project) ([#663](https://github.com/turbot/steampipe-plugin-gcp/pull/663)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)
 - [gcp_cloud_run_job](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_cloud_run_job) ([#673](https://github.com/turbot/steampipe-plugin-gcp/pull/673)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)
- Added `ignore_error_messages` config arg to provide users the ability to set a list of additional GCP error mesaages to ignore while running queries. For more information, please see [GCP plugin configuration](https://hub.steampipe.io/plugins/turbot/gcp#configuration).

_Bug fixes_

- Fixed the `unsupported region northamerica-south1` error while querying the `gcp_cloud_run_service` table. ([#707](https://github.com/turbot/steampipe-plugin-gcp/pull/707))
- Fixed the `akas` column of `gcp_redis_cluster` and `gcp_redis_instance` tables to correctly return data instead of `null`. ([#687](https://github.com/turbot/steampipe-plugin-gcp/pull/687))

## v1.2.0 [2024-11-04]

_What's new?_

- New tables added
  - [gcp_compute_instance_group_manager](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_instance_group_manager) ([#669](https://github.com/turbot/steampipe-plugin-gcp/pull/669)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)

_Enhancements_

- Added `labels` and `tags` columns to the `gcp_compute_global_forwarding_rule` table. ([#678](https://github.com/turbot/steampipe-plugin-gcp/pull/678)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)
- Added `database_installed_version` and `maintenance_version` columns to the `gcp_sql_database_instance` table. ([#677](https://github.com/turbot/steampipe-plugin-gcp/pull/677)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)

_Bug fixes_

- Fixed the `gcp_compute_instance_group` table to correctly return values for the `instances` column of regional instance groups. ([#670](https://github.com/turbot/steampipe-plugin-gcp/pull/670)) (Thanks [@pdecat](https://github.com/pdecat) for the contribution!)
- Fixed the `kubernetes_node_pool` table to correctly return data instead of an error for node pools with auto-pilot disabled. ([#668](https://github.com/turbot/steampipe-plugin-gcp/pull/668)) (Thanks [@multani](https://github.com/multani) for the contribution!)

## v1.0.0 [2024-10-22]

There are no significant changes in this plugin version; it has been released to align with [Steampipe's v1.0.0](https://steampipe.io/changelog/steampipe-cli-v1-0-0) release. This plugin adheres to [semantic versioning](https://semver.org/#semantic-versioning-specification-semver), ensuring backward compatibility within each major version.

_Breaking changes_

- The `credential_file` argument in the configuration file has now been removed. We recommend using the `credentials` argument instead, which can take the same file path as the `credential_file` argument. ([#672](https://github.com/turbot/steampipe-plugin-gcp/pull/672))

## v0.57.0 [2024-09-13]

_What's new?_

- New tables added
  - [gcp_composer_environment](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_composer_environment) ([#649](https://github.com/turbot/steampipe-plugin-gcp/pull/649))
  - [gcp_dataplex_lake](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_dataplex_lake) ([#642](https://github.com/turbot/steampipe-plugin-gcp/pull/642))
  - [gcp_dataplex_zone](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_dataplex_zone) ([#642](https://github.com/turbot/steampipe-plugin-gcp/pull/642))
  - [gcp_dataplex_task](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_dataplex_task) ([#642](https://github.com/turbot/steampipe-plugin-gcp/pull/642))
  - [gcp_dataplex_asset](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_dataplex_asset) ([#642](https://github.com/turbot/steampipe-plugin-gcp/pull/642))
  - [gcp_dataproc_metastore_service](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_dataproc_metastore_service) ([#654](https://github.com/turbot/steampipe-plugin-gcp/pull/654))
  - [gcp_vpc_access_connector](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_vpc_access_connector) ([#647](https://github.com/turbot/steampipe-plugin-gcp/pull/647))

_Enhancements_

- Added `connection_info` column to the `gcp_alloydb_instance` table. ([#651](https://github.com/turbot/steampipe-plugin-gcp/pull/651))

_Bug fixes_

- Removed the `name` column from the `gcp_bigquery_table` table since the API response does not include this field. ([#648](https://github.com/turbot/steampipe-plugin-gcp/pull/648))

_Dependencies_

- Recompiled plugin with Go version `1.22`. ([#635](https://github.com/turbot/steampipe-plugin-gcp/pull/635))
- Recompiled plugin with [steampipe-plugin-sdk v5.10.4](https://github.com/turbot/steampipe-plugin-sdk/blob/develop/CHANGELOG.md#v5104-2024-08-29) that fixes logging in the plugin export tool. ([#635](https://github.com/turbot/steampipe-plugin-gcp/pull/635))

## v0.56.0 [2024-07-31]

_Breaking changes_

- Removed the following columns in `gcp_cloudfunctions_function` table to align with the new API response structure: ([#612](https://github.com/turbot/steampipe-plugin-gcp/pull/612))
  - `environment_variables`
  - `source_upload_url`
  - `version_id`

_What's new?_

- Added the `impersonate_access_token` config argument to support plugin authentication by using a pre-generated temporary access token. ([#621](https://github.com/turbot/steampipe-plugin-gcp/pull/621))

_Enhancements_

- Added 17 new columns to the `gcp_cloudfunctions_function` table. ([#612](https://github.com/turbot/steampipe-plugin-gcp/pull/612))

_Bug fixes_

- Fixed the cache key issue in the `SecretManager` service client creation. ([#624](https://github.com/turbot/steampipe-plugin-gcp/pull/624))

## v0.55.0 [2024-07-22]

_Enhancements_

- Added column `create_time `to `gcp_sql_database_instance` table. ([#615](https://github.com/turbot/steampipe-plugin-gcp/pull/615))

_Bug fixes_

- Fixed the `gcp_alloydb_cluster` and `gcp_alloydb_instance` tables to correctly return values for `project` column instead of `null`. ([#617](https://github.com/turbot/steampipe-plugin-gcp/pull/617))

## v0.54.1 [2024-07-04]

_Bug fixes_

- Fixed the export tool of the plugin to return a non-zero error code instead of 0 whenever an error occurred. ([#79](https://github.com/turbot/steampipe-export/pull/79))

## v0.54.0 [2024-07-01]

_Bug fixes_

- Reverted the Export CLI behaviour to return `<nil>` for `null` values instead of `empty`. ([#77](https://github.com/turbot/steampipe-export/issues/77))

## v0.53.0 [2024-06-14]

_What's new?_

- New tables added
  - [gcp_alloydb_cluster](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_alloydb_cluster) ([#597](https://github.com/turbot/steampipe-plugin-gcp/pull/597))
  - [gcp_alloydb_instance](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_alloydb_cluster) ([#597](https://github.com/turbot/steampipe-plugin-gcp/pull/597))
  - [gcp_secret_manager_secret](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_secret_manager_secret) ([#593](https://github.com/turbot/steampipe-plugin-gcp/pull/593))

_Bug fixes_

- Fixed the `gcp_sql_database` table to skip SQL database instances that are not in the `running` state. ([#586](https://github.com/turbot/steampipe-plugin-gcp/pull/586))
- Fixed the `gcp_kubernetes_node_pool` table to skip clusters that have `autopilot` mode enabled. ([#591](https://github.com/turbot/steampipe-plugin-gcp/pull/591))

## v0.52.0 [2024-05-13]

_Enhancements_

- The `project` column has now been assigned as a connection key column across all the tables which facilitates more precise and efficient querying across multiple GCP projects. ([#564](https://github.com/turbot/steampipe-plugin-gcp/pull/564))
- The Plugin and the Steampipe Anywhere binaries are now built with the `netgo` package. ([#580](https://github.com/turbot/steampipe-plugin-gcp/pull/580))
- Added the `version` flag to the plugin's Export tool. ([#65](https://github.com/turbot/steampipe-export/pull/65))

_Bug fixes_

- Fixed the table `gcp_cloudfunctions_function` to list `gen2` cloud functions. ([#568](https://github.com/turbot/steampipe-plugin-gcp/pull/568)) (Thanks [@ashutoshmore658](https://github.com/ashutoshmore658) for the contribution!)

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.10.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v5100-2024-04-10) that adds support for connection key columns. ([#564](https://github.com/turbot/steampipe-plugin-gcp/pull/564))
- Recompiled plugin with [github.com/hashicorp/go-getter v1.7.4](https://github.com/hashicorp/go-getter/releases/tag/v1.7.4). ([#570](https://github.com/turbot/steampipe-plugin-gcp/pull/570))

## v0.51.1 [2024-04-12]

_Bug fixes_

- Fixed the plugin's Postgres FDW Extension crash [issue](https://github.com/turbot/steampipe-postgres-fdw/issues/434).

## v0.51.0 [2024-04-05]

_What's new?_

- New tables added
  - [gcp_tag_binding](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_tag_binding) ([#560](https://github.com/turbot/steampipe-plugin-gcp/pull/560))

## v0.50.0 [2024-03-29]

_What's new?_

- New tables added
  - [gcp_vertex_ai_model](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_vertex_ai_model) ([#555](https://github.com/turbot/steampipe-plugin-gcp/pull/555))

_Enhancements_

- Added support for `quota_project` config arg to provide users the ability to set the `project ID` used for billing and quota. ([#556](https://github.com/turbot/steampipe-plugin-gcp/pull/556))

_Bug fixes_

- Fixed the `retry_policy_maximum_backoff` and `retry_policy_minimum_backoff` columns of `gcp_pubsub_subscription` table to correctly return data. ([#552](https://github.com/turbot/steampipe-plugin-gcp/pull/552)) (Thanks to [@mvanholsteijn](https://github.com/mvanholsteijn) for the contribution!)

## v0.49.0 [2024-02-02]

_What's new?_

- New tables added
  - [gcp_app_engine_application](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_app_engine_application) ([#537](https://github.com/turbot/steampipe-plugin-gcp/pull/537))
  - [gcp_compute_machine_image](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_machine_image) ([#519](https://github.com/turbot/steampipe-plugin-gcp/pull/519))

## v0.48.0 [2024-01-22]

_What's new?_

- New tables added
  - [gcp_cloud_asset](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_cloud_asset) ([#532](https://github.com/turbot/steampipe-plugin-gcp/pull/532))

_Enhancements_

- Added column `iam_policy `to `gcp_cloud_run_service` table. ([#531](https://github.com/turbot/steampipe-plugin-gcp/pull/531))
- Optimized the `gcp_logging_log_entry` table result or result timing by applying a timestamp filter. ([#508](https://github.com/turbot/steampipe-plugin-gcp/pull/508))
- Added the `json_payload`, `proto_payload`, `metadata`, `resource`, `operation`, and `tags` columns to `gcp_logging_log_entry` table. ([#508](https://github.com/turbot/steampipe-plugin-gcp/pull/508))

_Bug fixes_

- Fixed the `addons_config`, `network_config` and `network_policy` column of `gcp_kubernetes_cluster` table to correctly return data instead of null. ([#530](https://github.com/turbot/steampipe-plugin-gcp/pull/530))
- Fixed the `end_time` column of the `gcp_sql_backup` table to return `null` instead of an error when end time is unavailable for a SQL backup. ([#534](https://github.com/turbot/steampipe-plugin-gcp/pull/534))
- Fixed the `enqueued_time`, `start_time` and `window_start_time` columns of the `gcp_sql_backup` table to return `null` instead of an error when timestamp is unavailable for a SQL backup. ([#536](https://github.com/turbot/steampipe-plugin-gcp/pull/536))

## v0.47.0 [2023-12-12]

_What's new?_

- The plugin can now be downloaded and used with the [Steampipe CLI](https://steampipe.io/install/steampipe.sh), as a [Postgres FDW](https://steampipe.io/install/postgres.sh), as a [SQLite extension](https://steampipe.io/install/sqlite.sh) and as a standalone [exporter](https://steampipe.io/install/export.sh).
- The table docs have been updated to provide corresponding example queries for Postgres FDW and SQLite extension.
- Docs license updated to match Steampipe [CC BY-NC-ND license](https://github.com/turbot/steampipe-plugin-gcp/blob/main/docs/LICENSE).

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.8.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v580-2023-12-11) that includes plugin server enacapsulation for in-process and GRPC usage, adding Steampipe Plugin SDK version to `_ctx` column, and fixing connection and potential divide-by-zero bugs. ([#521](https://github.com/turbot/steampipe-plugin-gcp/pull/521))

## v0.46.0 [2023-11-16]

_What's new?_

- New tables added
  - [gcp_aiplatform_endpoints](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_aiplatform_endpoints) ([#513](https://github.com/turbot/steampipe-plugin-gcp/pull/513))

_Bug fixes_

- Fixed the `retention_policy` column of `gcp_storage_bucket` table to correctly return data instead of null. ([#502](https://github.com/turbot/steampipe-plugin-gcp/pull/502))

## v0.45.0 [2023-10-04]

_What's new?_

- New tables added
  - [gcp_cloud_run_service](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_cloud_run_service) ([#497](https://github.com/turbot/steampipe-plugin-gcp/pull/497))
  - [gcp_artifact_registry_repository](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_artifact_registry_repository) ([#496](https://github.com/turbot/steampipe-plugin-gcp/pull/496))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v562-2023-10-03) which prevents nil pointer reference errors for implicit hydrate configs. ([#498](https://github.com/turbot/steampipe-plugin-gcp/pull/498))

## v0.44.0 [2023-10-02]

_Dependencies_

- Upgraded to [steampipe-plugin-sdk v5.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v561-2023-09-29) with support for rate limiters. ([#493](https://github.com/turbot/steampipe-plugin-gcp/pull/493))
- Recompiled plugin with Go version `1.21`. ([#493](https://github.com/turbot/steampipe-plugin-gcp/pull/493))

## v0.43.0 [2023-08-24]

_What's new?_

- New tables added
  - [gcp_redis_instance](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_redis_instance) ([#482](https://github.com/turbot/steampipe-plugin-gcp/pull/482))

## v0.42.1 [2023-08-17]

_Bug fixes_

- Fixed the `gcp_compute_machine_type` table to correctly return the supported machine types of all the zones. ([#480](https://github.com/turbot/steampipe-plugin-gcp/pull/480))

## v0.42.0 [2023-08-07]

_What's new?_

- New tables added
  - [gcp_cloud_identity_group](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_cloud_identity_group) ([#468](https://github.com/turbot/steampipe-plugin-gcp/pull/468))
  - [gcp_cloud_identity_group_membership](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_cloud_identity_group_membership) ([#468](https://github.com/turbot/steampipe-plugin-gcp/pull/468))

_Enhancements_

- Added `prefix` as an optional list key column in `gcp_storage_object` table. ([#475](https://github.com/turbot/steampipe-plugin-gcp/pull/475)) (Thanks to [@villers](https://github.com/villers) for the contribution!)

## v0.41.0 [2023-07-20]

_What's new?_

- New tables added
  - [gcp_logging_log_entry](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_logging_log_entry) ([#463](https://github.com/turbot/steampipe-plugin-gcp/pull/463))

## v0.40.0 [2023-06-29]

_What's new?_

- New tables added
  - [gcp_storage_object](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_storage_object) ([#460](https://github.com/turbot/steampipe-plugin-gcp/pull/460))

## v0.39.0 [2023-06-20]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.5.0](https://github.com/turbot/steampipe-plugin-sdk/blob/v5.5.0/CHANGELOG.md#v550-2023-06-16) which significantly reduces API calls and boosts query performance, resulting in faster data retrieval. ([#453](https://github.com/turbot/steampipe-plugin-gcp/pull/453))

## v0.38.0 [2023-06-15]

_What's new?_

- New tables added
  - [gcp_billing_account](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_billing_account) ([#443](https://github.com/turbot/steampipe-plugin-gcp/pull/443))
  - [gcp_billing_budget](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_billing_budget) ([#443](https://github.com/turbot/steampipe-plugin-gcp/pull/443))

## v0.37.0 [2023-05-19]

_Breaking changes_

- Removed columns `private_key_data`, `private_key_type` and `public_key_data` from `gcp_service_account_key` table as they consistently returned `null`. ([#449](https://github.com/turbot/steampipe-plugin-gcp/pull/449))

_Enhancements_

- Added columns `public_key_data_pem` and `public_key_data_raw` to `gcp_service_account_key` table. ([#449](https://github.com/turbot/steampipe-plugin-gcp/pull/449))

## v0.36.0 [2023-05-11]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.4.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v541-2023-05-05) which fixes increased plugin initialization time due to multiple connections causing the schema to be loaded repeatedly. ([#446](https://github.com/turbot/steampipe-plugin-gcp/pull/446))

## v0.35.0 [2023-04-05]

_What's new?_

- New tables added
  - [gcp_apikeys_key](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_apikeys_key) ([#440](https://github.com/turbot/steampipe-plugin-gcp/pull/440))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.3.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v530-2023-03-16) which adds go-getter support to dynamic tables. ([#444](https://github.com/turbot/steampipe-plugin-gcp/pull/444))

## v0.34.0 [2023-03-13]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v5.2.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v520-2023-03-02) which includes fixes for query cache pending item mechanism and aggregator connections not working for dynamic tables. ([#438](https://github.com/turbot/steampipe-plugin-gcp/pull/438))

## v0.33.1 [2023-03-06]

_Bug fixes_

- Fixed the plugin failing to detect when an active project is not set in the connection config, environment variables, or the gcloud CLI.

## v0.33.0 [2023-02-10]

_What's new?_

- Added `ignore_error_codes` config arg to provide users the ability to set a list of additional GCP error codes to ignore while running queries. For instance, to ignore some common access denied errors, which is helpful when running with limited permissions, set the argument `ignore_error_codes = ["401", 403"]`. With the addition of this config argument, it must be noted that the plugin will no longer return an empty row, when querying for services that have their APIs disabled, instead an error will be returned mentioning that the API has been disabled. To ignore this error, set the argument `ignore_error_codes = ["403"]`. For more information, please see [GCP plugin configuration](https://hub.steampipe.io/plugins/turbot/gcp#configuration). ([#432](https://github.com/turbot/steampipe-plugin-gcp/pull/432))

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.12](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v4112-2023-02-09) which fixes the query caching functionality. ([#434](https://github.com/turbot/steampipe-plugin-gcp/pull/434))

## v0.32.0 [2023-01-11]

_Breaking changes_

- Renamed column `name` to `key_name` in `gcp_kms_key_version` table to reflect that it displays the name of the KMS key rather than the name of the crypto key version. ([#424](https://github.com/turbot/steampipe-plugin-gcp/pull/424))

_Enhancements_

- Added column `id` to `gcp_kubernetes_cluster` table. ([#422](https://github.com/turbot/steampipe-plugin-gcp/pull/422))

_Bug fixes_

- Fixed the `id` column in `gcp_compute_instance_template`, `gcp_compute_network`, `gcp_compute_target_https_proxy` and `gcp_compute_target_pool` tables to be of `INT` datatype instead of `DOUBLE`. ([#426](https://github.com/turbot/steampipe-plugin-gcp/pull/426))
- Fixed the `title` column of `gcp_kms_key_version` table to return the KMS crypto key version appended to the KMS key name (for e.g. `my-key-name/1`) instead of returning only the key name (`my-key-name`). ([#424](https://github.com/turbot/steampipe-plugin-gcp/pull/424))

## v0.31.1 [2022-12-08]

_Bug fixes_

- Fixed `gcp_kubernetes_cluster` and `gcp_kubernetes_node_pool` tables always using application default credentials instead of user configured credentials. ([#420](https://github.com/turbot/steampipe-plugin-gcp/pull/420))

## v0.31.0 [2022-12-06]

_Enhancements_

- Updated `docs/index.md` to include better and more detailed multi-project configuration examples. ([#414](https://github.com/turbot/steampipe-plugin-gcp/pull/414))
- Updated the `Credentials` section of `docs/index.md` to include information about authenticating with Cloud Shell and Cloud Code. ([#395](https://github.com/turbot/steampipe-plugin-gcp/pull/395))

## v0.30.0 [2022-12-02]

_What's new?_

- New tables added
  - [gcp_compute_autoscaler](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_autoscaler) ([#412](https://github.com/turbot/steampipe-plugin-gcp/pull/412))
  - [gcp_compute_ha_vpn_gateway](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_ha_vpn_gateway) ([#413](https://github.com/turbot/steampipe-plugin-gcp/pull/413))
  - [gcp_compute_instance_group](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_instance_group) ([#401](https://github.com/turbot/steampipe-plugin-gcp/pull/401))

_Bug fixes_

- Fixed `gcp_compute_instance_metric_cpu_utilization_*` table docs to use the correct value of average in queries checking for CPU utilization. ([#416](https://github.com/turbot/steampipe-plugin-gcp/pull/416))

## v0.29.0 [2022-11-23]

_Dependencies_

- Recompiled plugin with [google-api-go-client v0.101.0](https://github.com/googleapis/google-api-go-client). ([#399](https://github.com/turbot/steampipe-plugin-gcp/pull/399))

## v0.28.1 [2022-11-11]

_Bug fixes_

- Fixed the `ip_address` column in `gcp_compute_forwarding_rule` and `gcp_compute_global_forwarding_rule` tables to be of `INET` data type instead of `IPADDR` data type. ([#409](https://github.com/turbot/steampipe-plugin-gcp/pull/409))
- Fixed the `gcp_cloudfunctions_function` table to correctly return data instead of an error when specifying the `name` column in the where clause. ([#403](https://github.com/turbot/steampipe-plugin-gcp/pull/403))
- Fixed the `Credentials` section of the `docs/index.md` to include the correct permission scope for querying the tables.

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.8](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v418-2022-09-08) which increases the default open file limit. ([#410](https://github.com/turbot/steampipe-plugin-gcp/pull/410))

## v0.28.0 [2022-10-24]

_Enhancements_

- Added `name` as optional ListConfig key column to `gcp_compute_disk_metric_*` and `gcp_compute_instance_metric_*` tables to improve query performance when filtering on disk or instance name respectively. ([#397](https://github.com/turbot/steampipe-plugin-gcp/pull/397)) (Thanks to [@JoshRosen](https://github.com/JoshRosen) for the improvements!)
- Updated examples for `project` config argument to use `PROJECT_ID` instead of `PROJECT_NAME` to better reflect the expected value. ([#398](https://github.com/turbot/steampipe-plugin-gcp/pull/398)) (Thanks to [@ashishsecdev](https://github.com/ashishsecdev) for the updates!)

_Bug fixes_

- Fixed `gcp_compute_disk` table queries panicking when querying the `iam_policy` column for a disk name that doesn't exist.

## v0.27.0 [2022-09-06]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v4.1.6](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v416-2022-09-02) which includes several caching and memory management improvements. ([#387](https://github.com/turbot/steampipe-plugin-gcp/pull/387))
- Recompiled plugin with Go version `1.19`. ([#383](https://github.com/turbot/steampipe-plugin-gcp/pull/383))

_Enhancements_

- Added column `self_link` to `gcp_dataproc_cluster table`. ([#386](https://github.com/turbot/steampipe-plugin-gcp/pull/386))

## v0.26.0 [2022-08-16]

_What's new?_

- New tables added
  - [gcp_dataproc_cluster](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_dataproc_cluster) ([#379](https://github.com/turbot/steampipe-plugin-gcp/pull/379))

_Enhancements_

- Added column `access_approval_settings` to `gcp_project` table. ([#381](https://github.com/turbot/steampipe-plugin-gcp/pull/381))
- Added column `essential_contacts` to `gcp_organization` table. ([#380](https://github.com/turbot/steampipe-plugin-gcp/pull/380))

## v0.25.0 [2022-07-13]

_Dependencies_

- Recompiled plugin with [steampipe-plugin-sdk v3.3.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v332--2022-07-11) which includes several caching fixes. ([#371](https://github.com/turbot/steampipe-plugin-gcp/pull/371))

## v0.24.0 [2022-07-01]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v3.3.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v331--2022-06-30). ([#368](https://github.com/turbot/steampipe-plugin-gcp/pull/368))

## v0.23.1 [2022-05-23]

_Bug fixes_

- Fixed the Slack community links in README and docs/index.md files. ([#363](https://github.com/turbot/steampipe-plugin-gcp/pull/363))

## v0.23.0 [2022-04-27]

_Enhancements_

- Added support for native Linux ARM and Mac M1 builds. ([#360](https://github.com/turbot/steampipe-plugin-gcp/pull/360))
- Recompiled plugin with [steampipe-plugin-sdk v3.1.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v310--2022-03-30) and Go version `1.18`. ([#356](https://github.com/turbot/steampipe-plugin-gcp/pull/356))

## v0.22.1 [2022-03-25]

_Bug fixes_

- Fixed the `gcp_compute_resource_policy` table to correctly return the compute resource policies instead of an error ([#352](https://github.com/turbot/steampipe-plugin-gcp/pull/352))

## v0.22.0 [2022-03-23]

_What's new?_

- New tables added
  - [gcp_kms_key_version](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_kms_key_version) ([#345](https://github.com/turbot/steampipe-plugin-gcp/pull/345))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v2.1.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v211--2022-03-10) ([#347](https://github.com/turbot/steampipe-plugin-gcp/pull/347))

## v0.21.0 [2022-03-17]

_Enhancements_

- Added the `autopilot_enabled` column to the `gcp_kubernetes_cluster` table ([#344](https://github.com/turbot/steampipe-plugin-gcp/pull/344))

_Bug fixes_

- The `gcp_project` table now correctly returns projects that have a project name different from their project ID

## v0.20.0 [2022-01-19]

_Enhancements_

- Added additional optional key quals, filter support, and context cancellation handling to `BigQuery`, `Cloud Functions`, `Compute Engine`, `DNS`, `KMS`, `Logging`, `Monitoring`, `Project`, `Pub/Sub` and `SQL` tables ([#337](https://github.com/turbot/steampipe-plugin-gcp/pull/337))
- Recompiled plugin with [steampipe-plugin-sdk v1.8.3](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v183--2021-12-23) ([#339](https://github.com/turbot/steampipe-plugin-gcp/pull/339))

## v0.19.0 [2021-11-23]

_Enhancements_

- Updated: Recompiled plugin with [steampipe-plugin-sdk v1.8.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v182--2021-11-22) and Go version 1.17 ([#332](https://github.com/turbot/steampipe-plugin-gcp/pull/332)) ([#335](https://github.com/turbot/steampipe-plugin-gcp/pull/335))

## v0.18.0 [2021-10-26]

_What's new?_

- Added: The `credentials` argument can now be specified in the configuration file to pass in either the path to or the contents of a service account key file in JSON format ([#318](https://github.com/turbot/steampipe-plugin-gcp/pull/318))

_Enhancements_

- Updated: Recompiled plugin with [steampipe-plugin-sdk v1.7.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md) ([#317](https://github.com/turbot/steampipe-plugin-gcp/pull/317))
- Updated: Add additional optional key quals to `gcp_iam_role`, `gcp_compute_image` and `gcp_compute_instance` tables ([#298](https://github.com/turbot/steampipe-plugin-gcp/pull/298))
- Updated: Add `WithCache` functionality for hydrate functions across all the tables ([#298](https://github.com/turbot/steampipe-plugin-gcp/pull/298))

_Deprecated_

- The `credential_file` argument in the configuration file is now deprecated and will be removed in the next major version. We recommend using the `credentials` argument instead, which can take the same file path as the `credential_file` argument. ([#318](https://github.com/turbot/steampipe-plugin-gcp/pull/318))

## v0.17.1 [2021-10-12]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.6.2](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v162--2021-10-08) ([#313](https://github.com/turbot/steampipe-plugin-gcp/pull/313))
- `credential_file` parameter in the configuration file now extends tildes to the current users home directory ([#312](https://github.com/turbot/steampipe-plugin-gcp/pull/312))

## v0.17.0 [2021-09-23]

_Enhancements_

- Added the `self_link` column to `gcp_bigtable_instance`, `gcp_cloudfunctions_function`, `gcp_monitoring_notification_channel`, `gcp_pubsub_snapshot`, `gcp_pubsub_subscription`, `gcp_pubsub_topic` and `gcp_project` tables ([#304](https://github.com/turbot/steampipe-plugin-gcp/pull/304))
- Added the `labels` column to `gcp_compute_instance_template` and `gcp_kms_key` tables ([#301](https://github.com/turbot/steampipe-plugin-gcp/pull/301))
- Recompiled plugin with [steampipe-plugin-sdk v1.6.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v161--2021-09-21) ([#309](https://github.com/turbot/steampipe-plugin-gcp/pull/309))

_Bug fixes_

- Column name `selfLink` has been correctly renamed to `self_link` in the `gcp_kubernetes_node_pool` table ([#308](https://github.com/turbot/steampipe-plugin-gcp/pull/308))

## v0.16.1 [2021-09-13]

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.5.1](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v151--2021-09-13) ([#303](https://github.com/turbot/steampipe-plugin-gcp/pull/303))

_Bug fixes_

- If the Google service API is disabled, then the corresponding table will now return zero records ([#287](https://github.com/turbot/steampipe-plugin-gcp/pull/287))

## v0.16.0 [2021-08-26]

_What's new?_

- New tables added
  - [gcp_kubernetes_cluster](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_kubernetes_cluster) ([#213](https://github.com/turbot/steampipe-plugin-gcp/pull/213))
  - [gcp_kubernetes_node_pool](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_kubernetes_node_pool) ([#214](https://github.com/turbot/steampipe-plugin-gcp/pull/214))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v1.5.0](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v150--2021-08-06) ([#290](https://github.com/turbot/steampipe-plugin-gcp/pull/290))

## v0.15.0 [2021-07-08]

_What's new?_

- New tables added
  - [gcp_sql_database_instance_metric_connections](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_sql_database_instance_metric_connections) ([#267](https://github.com/turbot/steampipe-plugin-gcp/pull/267))
  - [gcp_sql_database_instance_metric_connections_hourly](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_sql_database_instance_metric_connections_hourly) ([#268](https://github.com/turbot/steampipe-plugin-gcp/pull/268))

_Enhancements_

- Updated: Add columns `location` and `project` to all metric tables ([#279](https://github.com/turbot/steampipe-plugin-gcp/pull/279))

_Bug fixes_

- Fixed: `timestamp` column type in all metric tables is now `TIMESTAMP` instead of `STRING` ([#277](https://github.com/turbot/steampipe-plugin-gcp/pull/277))
- Fixed: Typo in table name and invalid key column definition for `gcp_dns_record_set table`

## v0.14.0 [2021-07-01]

_What's new?_

- New tables added
  - [gcp_bigquery_job](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_bigquery_job) ([#247](https://github.com/turbot/steampipe-plugin-gcp/pull/247))
  - [gcp_compute_disk_metric_read_ops](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_disk_metric_read_ops) ([#262](https://github.com/turbot/steampipe-plugin-gcp/pull/262))
  - [gcp_compute_disk_metric_read_ops_daily](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_disk_metric_read_ops_daily) ([#264](https://github.com/turbot/steampipe-plugin-gcp/pull/264))
  - [gcp_compute_disk_metric_read_ops_hourly](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_disk_metric_read_ops_hourly) ([#260](https://github.com/turbot/steampipe-plugin-gcp/pull/260))
  - [gcp_compute_disk_metric_write_ops](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_disk_metric_write_ops) ([#272](https://github.com/turbot/steampipe-plugin-gcp/pull/272))
  - [gcp_compute_disk_metric_write_ops_daily](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_disk_metric_write_ops_daily) ([#263](https://github.com/turbot/steampipe-plugin-gcp/pull/263))
  - [gcp_compute_disk_metric_write_ops_hourly](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_disk_metric_write_ops_hourly) ([#271](https://github.com/turbot/steampipe-plugin-gcp/pull/271))
  - [gcp_compute_instance_metric_cpu_utilization](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_instance_metric_cpu_utilization) ([#254](https://github.com/turbot/steampipe-plugin-gcp/pull/254))
  - [gcp_compute_instance_metric_cpu_utilization_daily](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_instance_metric_cpu_utilization_daily) ([#256](https://github.com/turbot/steampipe-plugin-gcp/pull/256))
  - [gcp_compute_instance_metric_cpu_utilization_hourly](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_instance_metric_cpu_utilization_hourly) ([#255](https://github.com/turbot/steampipe-plugin-gcp/pull/255))
  - [gcp_compute_machine_type](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_machine_type) ([#238](https://github.com/turbot/steampipe-plugin-gcp/pull/238))
  - [gcp_compute_resource_policy](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_resource_policy) ([#237](https://github.com/turbot/steampipe-plugin-gcp/pull/237))
  - [gcp_dns_record_set](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_dns_record_set) ([#212](https://github.com/turbot/steampipe-plugin-gcp/pull/212))
  - [gcp_logging_bucket](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_logging_bucket) ([#232](https://github.com/turbot/steampipe-plugin-gcp/pull/232))
  - [gcp_organization](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_organization) ([#250](https://github.com/turbot/steampipe-plugin-gcp/pull/250))
  - [gcp_project_organization_policy](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_project_organization_policy) ([#204](https://github.com/turbot/steampipe-plugin-gcp/pull/204))
  - [gcp_sql_database_instance_metric_connections_daily](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_sql_database_instance_metric_connections_daily) ([#240](https://github.com/turbot/steampipe-plugin-gcp/pull/240))
  - [gcp_sql_database_instance_metric_cpu_utilization](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_sql_database_instance_metric_cpu_utilization) ([#259](https://github.com/turbot/steampipe-plugin-gcp/pull/259))
  - [gcp_sql_database_instance_metric_cpu_utilization_daily](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_sql_database_instance_metric_cpu_utilization_daily) ([#251](https://github.com/turbot/steampipe-plugin-gcp/pull/251))
  - [gcp_sql_database_instance_metric_cpu_utilization_hourly](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_sql_database_instance_metric_cpu_utilization_hourly) ([#261](https://github.com/turbot/steampipe-plugin-gcp/pull/261))

_Enhancements_

- Updated: Add support for service account impersonation credentials ([#249](https://github.com/turbot/steampipe-plugin-gcp/pull/249))

_Bug fixes_

- Fixed: List and get calls should not error in `gcp_sql_database` table for stopped databases ([#244](https://github.com/turbot/steampipe-plugin-gcp/pull/244))

## v0.13.1 [2021-06-17]

_Bug fixes_

- Fixed: `id` column type in `gcp_compute_disk` table ([#228](https://github.com/turbot/steampipe-plugin-gcp/pull/228))
- Fixed: Example query in `gcp_compute_address` table doc ([#224](https://github.com/turbot/steampipe-plugin-gcp/pull/224))

## v0.13.0 [2021-05-27]

_What's new?_

- Updated plugin license to Apache 2.0 per [turbot/steampipe#488](https://github.com/turbot/steampipe/issues/488)
- New tables added
  - [gcp_kms_key](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_kms_key) ([#215](https://github.com/turbot/steampipe-plugin-gcp/pull/215))

## v0.12.0 [2021-05-20]

_What's new?_

- New tables added
  - [gcp_compute_project_metadata](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_project_metadata) ([#192](https://github.com/turbot/steampipe-plugin-gcp/pull/192))
  - [gcp_dns_policy](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_dns_policy) ([#208](https://github.com/turbot/steampipe-plugin-gcp/pull/208))
  - [gcp_project](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_project) ([#205](https://github.com/turbot/steampipe-plugin-gcp/pull/205))

## v0.11.0 [2021-05-13]

_What's new?_

- New tables added
  - [gcp_bigquery_table](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_bigquery_table) ([#198](https://github.com/turbot/steampipe-plugin-gcp/pull/198))

_Enhancements_

- Updated: README.md and docs/index.md now contain links to our Slack community ([#202](https://github.com/turbot/steampipe-plugin-gcp/pull/202))
- Updated: Add `self_link` column to `gcp_dns_managed_zone` table ([#195](https://github.com/turbot/steampipe-plugin-gcp/pull/195))
- Updated: Add `self_link` column to `gcp_logging_sink` table ([#193](https://github.com/turbot/steampipe-plugin-gcp/pull/193))
- Updated: Bump lodash from 4.17.20 to 4.17.21 in /gcp-test ([#197](https://github.com/turbot/steampipe-plugin-gcp/pull/197))

## v0.10.0 [2021-05-06]

_What's new?_

_Enhancements_

- Updated: Add `disk_encryption_key` and `disk_encryption_key_type` columns to `gcp_compute_disk` table ([#185](https://github.com/turbot/steampipe-plugin-gcp/pull/185))
- Updated: Remove `disk_encryption_kms_key` column from `gcp_compute_disk` table ([#185](https://github.com/turbot/steampipe-plugin-gcp/pull/185))
- Updated: Add `metric_descriptor_type` column to `gcp_logging_metric` table ([#182](https://github.com/turbot/steampipe-plugin-gcp/pull/182))

## v0.9.0 [2021-04-29]

_What's new?_

- New tables added
  - [gcp_compute_target_ssl_proxy](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_target_ssl_proxy) ([#156](https://github.com/turbot/steampipe-plugin-gcp/pull/156))

## v0.8.0 [2021-04-22]

_What's new?_

- New tables added
  - [gcp_kms_key_ring](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_kms_key_ring) ([#171](https://github.com/turbot/steampipe-plugin-gcp/pull/171))

## v0.7.0 [2021-04-15]

_What's new?_

- New tables added
  - [gcp_compute_ssl_policy](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_ssl_policy) ([#154](https://github.com/turbot/steampipe-plugin-gcp/pull/154))
  - [gcp_monitoring_alert_policy](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_monitoring_alert_policy) ([#168](https://github.com/turbot/steampipe-plugin-gcp/pull/168))

_Bug fixes_

- Fixed: Query example in `gcp_dns_managed_zone` table docs ([#172](https://github.com/turbot/steampipe-plugin-gcp/pull/172))

## v0.6.0 [2021-04-08]

_What's new?_

- New tables added
  - [gcp_compute_target_https_proxy](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_compute_target_https_proxy) ([#151](https://github.com/turbot/steampipe-plugin-gcp/pull/151))

_Enhancements_

- Updated: `gcp_sql_backup`, `gcp_sql_database`, `gcp_sql_database_instance` tables now use the `sqladmin` package instead of the `sql` package ([#161](https://github.com/turbot/steampipe-plugin-gcp/pull/161))
- Recompiled plugin with [steampipe-plugin-sdk v0.2.7](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v027-2021-03-31)

_Bug fixes_

- Fixed: Backup configuration columns now return the correct data in the `gcp_sql_database_instance` table ([#166](https://github.com/turbot/steampipe-plugin-gcp/pull/166))
- Fixed: Removed unused `root_password` column in the `gcp_sql_database_instance` table ([#166](https://github.com/turbot/steampipe-plugin-gcp/pull/166))

## v0.5.1 [2021-04-02]

_Bug fixes_

- Fixed: `Table definitions & examples` link now points to the correct location ([#163](https://github.com/turbot/steampipe-plugin-gcp/pull/163))

## v0.5.0 [2021-04-01]

_What's new?_

- New tables added
  - [gcp_bigquery_dataset](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_bigquery_dataset) ([#146](https://github.com/turbot/steampipe-plugin-gcp/pull/146))
  - [gcp_dns_managed_zone](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_dns_managed_zone) ([#147](https://github.com/turbot/steampipe-plugin-gcp/pull/147))

## v0.4.0 [2021-03-25]

_What's new?_

- New tables added
  - [gcp_bigtable_instance](https://hub.steampipe.io/plugins/turbot/gcp/tables/gcp_bigtable_instance) ([#90](https://github.com/turbot/steampipe-plugin-gcp/pull/90))

_Enhancements_

- Recompiled plugin with [steampipe-plugin-sdk v0.2.6](https://github.com/turbot/steampipe-plugin-sdk/blob/main/CHANGELOG.md#v026-2021-03-18)

## v0.3.1 [2021-03-11]

_Bug fixes_

- Removed use of deprecated `ItemFromKey` function from all tables

## v0.3.0 [2021-03-04]

_What's new?_

- New tables added
  - gcp_compute_region
  - gcp_compute_zone

## v0.2.1 [2021-03-02]

_Bug fixes_

- Recompiled plugin with latest [steampipe-plugin-sdk](https://github.com/turbot/steampipe-plugin-sdk) to resolve issue:
  - Fix tables failing with error similar to `Error: pq: rpc error: code = Internal desc = get hydrate function getIamRole failed with panic interface conversion: interface {} is nil, not *gcp.roleInfo`([#129](https://github.com/turbot/steampipe-plugin-gcp/issues/129)).

## v0.2.0 [2021-02-25]

_What's new?_

- New tables added
  - gcp_sql_backup
  - gcp_sql_database
  - gcp_sql_database_instance

_Bug fixes_

- Updated `gcp_compute_instance` table `network_tags` field to display value correctly ([#114](https://github.com/turbot/steampipe-plugin-gcp/pull/114))
- Recompiled plugin with latest [steampipe-plugin-sdk](https://github.com/turbot/steampipe-plugin-sdk) to resolve SDK issues:
  - Fix error for missing required quals [#40](https://github.com/turbot/steampipe-plugin-sdk/issues/42).
  - Queries fail with error socket: too many open files [#190](https://github.com/turbot/steampipe/issues/190)

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

#!/bin/bash
 RED="\e[31m"
 GREEN="\e[32m"
 BLACK="\e[30m"
 BLUE="\e[34m"
 BOLDGREEN="\e[1;${GREEN}"
 PINK="\e[38;5;198m"
 ENDCOLOR="\e[0m"


# Define your function here
run_test () {
   echo -e "${PINK}Running $1 ${ENDCOLOR}"
 if ! ./tint.js $1 >> output.txt
   then
    echo -e "${RED}Failed -> $1 ${ENDCOLOR}"
    echo $1 >> failed_tests.txt
  else
    echo -e "${BLUE}Passed -> $1 ${ENDCOLOR}"
    echo $1 >> passed_tests.txt
   fi
 }

 # output.txt - store output of each test
 # failed_tests.txt - names of failed test
 # passed_tests.txt names of passed test

 # removes files from previous test
# rm -rf output.txt failed_tests.txt passed_tests.txt
 date >> output.txt
 date >> failed_tests.txt
 date >> passed_tests.txt

run_test gcp_bigquery_dataset
run_test gcp_bigquery_job
run_test gcp_bigquery_table
run_test gcp_bigtable_instance
run_test gcp_compute_address
run_test gcp_compute_autoscaler
run_test gcp_compute_backend_bucket
run_test gcp_compute_backend_service
run_test gcp_compute_disk
run_test gcp_compute_firewall
run_test gcp_compute_forwarding_rule
run_test gcp_compute_global_address
run_test gcp_compute_global_forwarding_rule
run_test gcp_compute_ha_vpn_gateway
run_test gcp_compute_image
run_test gcp_compute_instance
run_test gcp_compute_instance_group
run_test gcp_compute_instance_template
run_test gcp_compute_machine_image
run_test gcp_compute_machine_type
run_test gcp_compute_network
run_test gcp_compute_node_group
run_test gcp_compute_node_template
run_test gcp_compute_project_metadata
run_test gcp_compute_resource_policy
run_test gcp_compute_route
run_test gcp_compute_router
run_test gcp_compute_snapshot
run_test gcp_compute_ssl_policy
run_test gcp_compute_subnetwork
run_test gcp_compute_target_https_proxy
run_test gcp_compute_target_pool
run_test gcp_compute_target_ssl_proxy
run_test gcp_compute_target_vpn_gateway
run_test gcp_compute_url_map
run_test gcp_compute_vpn_tunnel
run_test gcp_dataproc_cluster
run_test gcp_dns_managed_zone
run_test gcp_dns_policy
run_test gcp_dns_record_set
run_test gcp_iam_role
run_test gcp_kms_key
run_test gcp_kms_key_ring
run_test gcp_kms_key_version
run_test gcp_kubernetes_cluster
run_test gcp_kubernetes_node_pool
run_test gcp_logging_bucket
run_test gcp_logging_exclusion
run_test gcp_logging_metric
run_test gcp_logging_sink
run_test gcp_monitoring_alert_policy
run_test gcp_monitoring_group
run_test gcp_monitoring_notification_channel
run_test gcp_organization
run_test gcp_project
run_test gcp_project_organization_policy
run_test gcp_pubsub_subscription
run_test gcp_pubsub_topic
run_test gcp_service_account
run_test gcp_service_account_key
run_test gcp_sql_backup
run_test gcp_sql_database
run_test gcp_sql_database_instance
run_test gcp_storage_bucket
run_test gcp_tag_binding

date >> output.txt
date >> failed_tests.txt
date >> passed_tests.txt
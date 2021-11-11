
#!/bin/bash
 RED="\e[31m"
 GREEN="\e[32m"
 BLACK="\e[30m"
 BOLDGREEN="\e[1;${GREEN}"
 ENDCOLOR="\e[0m"


# Define your function here
run_test () {
   echo -e "${BLACK}Running $1 ${ENDCOLOR}"
 if ! ./tint.js $1 >> output.txt
   then
    echo -e "${RED}Failed -> $1 ${ENDCOLOR}"
    echo $1 >> failed_tests.txt
  else
    echo -e "${BOLDGREEN}Passed -> $1 ${ENDCOLOR}"
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

run_test bigquery_dataset
run_test bigquery_job
run_test bigquery_table
run_test bigtable_instance
run_test compute_address
run_test compute_backend_bucket
run_test compute_backend_service
run_test compute_disk
run_test compute_firewall
run_test compute_forwarding_rule
run_test compute_global_address
run_test compute_global_forwarding_rule
run_test compute_image
run_test compute_instance
run_test compute_instance_template
run_test compute_machine_type
run_test compute_network
run_test compute_node_group
run_test compute_node_template
run_test compute_project_metadata
run_test compute_resource_policy
run_test compute_route
run_test compute_router
run_test compute_snapshot
run_test compute_ssl_policy
run_test compute_subnetwork
run_test compute_target_https_proxy
run_test compute_target_pool
run_test compute_target_ssl_proxy
run_test compute_target_vpn_gateway
run_test compute_url_map
run_test compute_vpn_tunnel
run_test dns_managed_zone
run_test dns_policy
run_test dns_record_set
run_test iam_role
run_test kms_key
run_test kms_key_ring
run_test kubernetes_cluster
run_test kubernetes_node_pool
run_test logging_bucket
run_test logging_exclusion
run_test logging_metric
run_test logging_sink
run_test monitoring_alert_policy
run_test monitoring_group
run_test monitoring_notification_channel
run_test organization
run_test project
run_test project_organization_policy
run_test pubsub_subscription
run_test pubsub_topic
run_test service_account
run_test service_account_key
run_test sql_backup
run_test sql_database
run_test sql_database_instance
run_test storage_bucket

date >> output.txt
date >> failed_tests.txt
date >> passed_tests.txt

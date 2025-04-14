---
title: "Steampipe Table: gcp_compute_machine_image - Query Google Cloud Platform Compute Machine Image using SQL"
description: "Allows users to query Compute Machine Images in Google Cloud Platform, providing detailed information about available machine images and their specifications."
folder: "Compute"
---

# Table: gcp_compute_machine_image - Query Google Cloud Platform Compute Machine Image using SQL

A machine image is a Compute Engine resource that stores all the configuration, metadata, permissions, and data from multiple disks of a virtual machine (VM) instance. You can use a machine image in many system maintenance, backup and recovery, and instance cloning scenarios.

## Table Usage Guide

The `gcp_compute_machine_image` table provides insights into the available machine images within Google Cloud Platform's Compute Engine. As a cloud architect or DevOps engineer, you can explore machine image-specific details through this table, kind, source instance, instance properties, image status, image storage, and associated metadata. Utilize it to understand the specifications of each machine image, aiding in the selection of the most suitable machine image for your applications based on performance requirements and cost efficiency.

## Examples

### Basic info
Assess the elements within your Google Cloud Platform to understand the capacity and capabilities of each machine image. This can help to get the metadata about the compute images.

```sql+postgres
select
  name,
  id,
  description,
  creation_timestamp,
  guest_flush,
  source_instance
from
  gcp_compute_machine_image;
```

```sql+sqlite
select
  name,
  id,
  description,
  creation_timestamp,
  guest_flush,
  source_instance
from
  gcp_compute_machine_image;
```

### List machine images that are available
Ensures that only machine images that are ready for deployment or use are considered, which is critical for operational stability and reliability. Useful in automated scripts or applications where only machine images in a 'READY' state should be utilized. Helps in maintaining a clean and efficient image repository by focusing on images that are fully prepared and excluding those that are still in preparation or have been deprecated.

```sql+postgres
select
  name,
  id,
  description,
  creation_timestamp,
  status
from
  gcp_compute_machine_image
where
  status = 'READY';
```

```sql+sqlite
select
  name,
  id,
  description,
  creation_timestamp,
  status
from
  gcp_compute_machine_image
where
  status = 'READY';
```

### List the top 5 machine images that consume highest storage
This query is particularly useful in cloud infrastructure management and optimization, where understanding and managing storage utilization is a key concern. It helps administrators and users quickly identify the most space-efficient machine images available in their GCP environment.

```sql+postgres
select
  name,
  id,
  self_link,
  status,
  total_storage_bytes
from
  gcp_compute_machine_image
order by
  total_storage_bytes asc
limit 5;
```

```sql+sqlite
select
  name,
  id,
  self_link,
  status,
  total_storage_bytes
from
  gcp_compute_machine_image
order by
  total_storage_bytes asc
limit 5;
```

### Get instance properties of the machine images
Useful for analyzing the detailed configurations of machine images, including hardware features, network settings, and security configurations. Assists in planning and optimizing cloud infrastructure based on the capabilities and configurations of available machine images.

```sql+postgres
select
  name,
  id,
  instance_properties -> 'advancedMachineFeatures' as advanced_machine_features,
  instance_properties ->> 'canIpForward' as can_ip_forward,
  instance_properties -> 'confidentialInstanceConfig' as confidential_instance_config,
  instance_properties ->> 'description' as description,
  instance_properties -> 'disks' as disks,
  instance_properties -> 'guestAccelerators' as guest_accelerators,
  instance_properties ->> 'keyRevocationActionType' as key_revocation_action_type,
  instance_properties -> 'labels' as labels,
  instance_properties ->> 'machineType' as machine_type,
  instance_properties -> 'metadata' as metadata,
  instance_properties -> 'minCpuPlatform' as min_cpu_platform,
  instance_properties -> 'networkInterfaces' as network_interfaces,
  instance_properties -> 'networkPerformanceConfig' as network_performance_config,
  instance_properties -> 'privateIpv6GoogleAccess' as private_ipv6_google_access,
  instance_properties ->> 'reservationAffinity' as reservation_affinity,
  instance_properties -> 'resourceManagerTags' as resource_manager_tags,
  instance_properties -> 'resourcePolicies' as resource_policies,
  instance_properties -> 'scheduling' as scheduling,
  instance_properties -> 'serviceAccounts' as service_accounts,
  instance_properties -> 'shieldedInstanceConfig' as shielded_instance_config,
  instance_properties -> 'tags' as tags
from
  gcp_compute_machine_image;
```

```sql+sqlite
select
  name,
  id,
  json_extract(instance_properties, '$.advancedMachineFeatures') as advanced_machine_features,
  json_extract(instance_properties, '$.canIpForward') as can_ip_forward,
  json_extract(instance_properties, '$.confidentialInstanceConfig') as confidential_instance_config,
  json_extract(instance_properties, '$.description') as description,
  json_extract(instance_properties, '$.disks') as disks,
  json_extract(instance_properties, '$.guestAccelerators') as guest_accelerators,
  json_extract(instance_properties, '$.keyRevocationActionType') as key_revocation_action_type,
  json_extract(instance_properties, '$.labels') as labels,
  json_extract(instance_properties, '$.machineType') as machine_type,
  json_extract(instance_properties, '$.metadata') as metadata,
  json_extract(instance_properties, '$.minCpuPlatform') as min_cpu_platform,
  json_extract(instance_properties, '$.networkInterfaces') as network_interfaces,
  json_extract(instance_properties, '$.networkPerformanceConfig') as network_performance_config,
  json_extract(instance_properties, '$.privateIpv6GoogleAccess') as private_ipv6_google_access,
  json_extract(instance_properties, '$.reservationAffinity') as reservation_affinity,
  json_extract(instance_properties, '$.resourceManagerTags') as resource_manager_tags,
  json_extract(instance_properties, '$.resourcePolicies') as resource_policies,
  json_extract(instance_properties, '$.scheduling') as scheduling,
  json_extract(instance_properties, '$.serviceAccounts') as service_accounts,
  json_extract(instance_properties, '$.shieldedInstanceConfig') as shielded_instance_config,
  json_extract(instance_properties, '$.tags') as tags
from
  gcp_compute_machine_image;
```

### Get encryption details of the machine image
Understanding the encryption methods and keys used for each machine image is vital for security and compliance. It helps ensure that sensitive data is properly protected and that the encryption methods meet required standards. The query aids in auditing the encryption practices and managing the encryption keys across different machine images. It's particularly useful in environments with strict data protection policies.

```sql+postgres
select
  name,
  machine_image_encryption_key ->> 'KmsKeyName' as kms_key_name,
  machine_image_encryption_key ->> 'KmsKeyServiceAccount' as kms_key_service_account,
  machine_image_encryption_key ->> 'RawKey' as raw_key,
  machine_image_encryption_key ->> 'RsaEncryptedKey' as rsa_encrypted_key,
  machine_image_encryption_key ->> 'Sha256' as sha256
from
  gcp_compute_machine_image;
```

```sql+sqlite
select
  name,
  json_extract(machine_image_encryption_key, '$.KmsKeyName') as kms_key_name,
  json_extract(machine_image_encryption_key, '$.KmsKeyServiceAccount') as kms_key_service_account,
  json_extract(machine_image_encryption_key, '$.RawKey') as raw_key,
  json_extract(machine_image_encryption_key, '$.RsaEncryptedKey') as rsa_encrypted_key,
  json_extract(machine_image_encryption_key, '$.Sha256') as sha256
from
  gcp_compute_machine_image;
```

### Get the machine type details for the machine images
Analyzing the memory, CPU, and disk capabilities of machine types can inform decisions about image deployment based on performance needs. Knowing the deprecation status and creation timestamp of machine types helps in compliance and migration planning.

```sql+postgres
select
  i.name as image_name,
  i.id image_id,
  i.instance_properties ->> 'machineType' as machine_type,
  t.creation_timestamp as machine_type_creation_timestamp,
  t.memory_mb as machine_type_memory_mb,
  t.maximum_persistent_disks as machine_type_maximum_persistent_disks,
  t.is_shared_cpu as machine_type_is_shared_cpu,
  t.zone as machine_type_zone,
  t.deprecated as machine_type_deprecated
from
  gcp_compute_machine_image as i,
  gcp_compute_machine_type as t
where
  t.name = (i.instance_properties ->> 'machineType') and t.zone = split_part(i.source_instance, '/', 9);
```

```sql+sqlite
select
  i.name as image_name,
  i.id as image_id,
  json_extract(i.instance_properties, '$.machineType') as machine_type,
  t.creation_timestamp as machine_type_creation_timestamp,
  t.memory_mb as machine_type_memory_mb,
  t.maximum_persistent_disks as machine_type_maximum_persistent_disks,
  t.is_shared_cpu as machine_type_is_shared_cpu,
  t.zone as machine_type_zone,
  t.deprecated as machine_type_deprecated
from
  gcp_compute_machine_image as i,
  gcp_compute_machine_type as t
where
  t.name = json_extract(i.instance_properties, '$.machineType')
  and t.zone = substr(i.source_instance, instr(i.source_instance, '/', -1) + 1);
```
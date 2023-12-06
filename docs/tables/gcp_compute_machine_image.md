# Table: gcp_compute_machine_image

A machine image is a Compute Engine resource that stores all the configuration, metadata, permissions, and data from multiple disks of a virtual machine (VM) instance. You can use a machine image in many system maintenance, backup and recovery, and instance cloning scenarios.

## Examples

### Basic info

```sql
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

```sql
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

```sql
select
  name,
  id,
  self_link,
  status,
  total_storage_bytes
from
  gcp_compute_machine_image
order by
  total_storage_bytes asc;
```

### Get instance properties of the machine images

```sql
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

### Get encryption details of the machine image

```sql
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

### Get the machine type details for the machine images

```sql
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
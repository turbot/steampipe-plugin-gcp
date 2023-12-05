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
  instance_properties -> 'AdvancedMachineFeatures' as advanced_machine_features,
  instance_properties ->> 'CanIpForward' as can_ip_forward,
  instance_properties -> 'ConfidentialInstanceConfig' as confidential_instance_config,
  instance_properties ->> 'Description' as description,
  instance_properties -> 'Disks' as disks,
  instance_properties -> 'GuestAccelerators' as guest_accelerators,
  instance_properties ->> 'KeyRevocationActionType' as key_revocation_action_type,
  instance_properties -> 'Labels' as labels,
  instance_properties ->> 'MachineType' as machine_type,
  instance_properties -> 'Metadata' as metadata,
  instance_properties -> 'MinCpuPlatform' as min_cpu_platform,
  instance_properties -> 'NetworkInterfaces' as network_interfaces,
  instance_properties -> 'NetworkPerformanceConfig' as network_performance_config,
  instance_properties -> 'PrivateIpv6GoogleAccess' as private_ipv6_google_access,
  instance_properties ->> 'ReservationAffinity' as reservation_affinity,
  instance_properties -> 'ResourceManagerTags' as resource_manager_tags,
  instance_properties -> 'ResourcePolicies' as resource_policies,
  instance_properties -> 'Scheduling' as scheduling,
  instance_properties -> 'ServiceAccounts' as service_accounts,
  instance_properties -> 'ShieldedInstanceConfig' as shielded_instance_config,
  instance_properties -> 'Tags' as tags
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
  machine_image_encryption_key ->> 'Sha256' as sha256,
from
  gcp_compute_machine_image;
```
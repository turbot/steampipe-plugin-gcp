# Table:  gcp_compute_instance_template

An instance template is a resource that use to create virtual machine (VM) instances and managed instance groups (MIGs). Instance templates define the machine type, boot disk image or container image, labels, and other instance properties.

## Examples

### List of c2-standard-4 machine type instance template

```sql
select
  name,
  id,
  instance_machine_type
from
  gcp_compute_instance_template
where
  instance_machine_type = 'c2-standard-4';
```


### Boot Disk info of each instance template
```sql
select
  name,
  id,
  disk ->> 'deviceName' as disk_device_name,
  disk -> 'initializeParams' ->> 'diskType' as disk_type,
  disk -> 'initializeParams' ->> 'diskSizeGb' as disk_size_gb,
  split_part(
    disk -> 'initializeParams' ->> 'sourceImage',
    '/',
    5
  ) as source_image,
  disk ->> 'mode' as mode
from
  gcp_compute_instance_template,
  jsonb_array_elements(instance_disks) as disk;
```


### List of SPECIFIC_RESERVATION Instance type instance template
```sql
select
  name,
  id,
  instance_reservation_affinity ->> 'consumeReservationType' as consume_reservation_type
from
  gcp_compute_instance_template
where
  instance_reservation_affinity ->> 'consumeReservationType' = 'SPECIFIC_RESERVATION';
```


### Network interface info of each instance template

```sql
select
  name,
  id,
  i ->> 'name' as name,
  split_part(i ->> 'network', '/', 10) as network_name,
  p ->> 'name' as access_config_name,
  p ->> 'networkTier' as access_config_network_tier,
  p ->> 'type' as access_config_type
from
  gcp_compute_instance_template,
  jsonb_array_elements(instance_network_interfaces) as i,
  jsonb_array_elements(i -> 'accessConfigs') as p;
```


### List of instance templates where instance_can_ip_forward is true

```sql
select
  name,
  id,
  instance_can_ip_forward
from
  gcp_compute_instance_template
where
  instance_can_ip_forward;
```
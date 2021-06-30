# Table: gcp_compute_machine_type

A machine type is a set of virtualized hardware resources available to a virtual machine (VM) instance, including the system memory size, virtual CPU (vCPU) count, and persistent disk limits.

In Compute Engine, machine types are grouped and curated by families for different workloads. Compute Engine offers general-purpose, memory-optimized, compute-optimized, and accelerated-optimized families.

## Examples

### Basic info

```sql
select
  name,
  id,
  description,
  guest_cpus,
  maximum_persistent_disks,
  maximum_persistent_disks_size_gb
from
  gcp_compute_machine_type;
```


### List machine types with more than 48 cores

```sql
select
  name,
  id,
  description,
  guest_cpus
from
  gcp_compute_machine_type
where
  guest_cpus >= 48;
```


### List machine types with shared CPUs

```sql
select
  name,
  id,
  is_shared_cpu
from
  gcp_compute_machine_type
where
  is_shared_cpu;
```


### Get accelerator configurations assigned to each machine type

```sql
select
  name,
  id,
  a -> 'guestAcceleratorCount' as guest_accelerator_count,
  a ->> 'guestAcceleratorType' as guest_accelerator_type
from
  gcp_compute_machine_type,
  jsonb_array_elements(accelerators) as a;
```

# Table: gcp_compute_zone

Compute Engine resources are hosted in multiple locations worldwide. These locations are composed of regions and zones. Resources that live in a zone, such as virtual machine instances or zonal persistent disks, are referred to as zonal resources.

## Examples

### Compute zone basic info

```sql
select
  name,
  id,
  region_name,
  status
from
  gcp_compute_zone;
```


### Get the available cpu platforms in each zone

```sql
select
  name,
  available_cpu_platforms
from
  gcp_compute_zone;
```


### Get the zones which are down

```sql
select
  name,
  id,
  region_name,
  status
from
  gcp_compute_zone
where
  status = 'DOWN';
```
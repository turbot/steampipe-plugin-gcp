# Table: gcp_compute_region

Compute Engine resources are hosted in multiple locations worldwide. These locations are composed of regions and zones. A region is a specific geographical location where users can host their resources.

## Examples

### List of compute regions which are down

```sql
select
  name,
  id,
  status
from
  gcp_compute_region
where
  status = 'DOWN';
```


### Get the quota info of each region

```sql
select
  name,
  q -> 'metric' as quota_metric,
  q -> 'limit' as quota_limit
from
  gcp_compute_region,
  jsonb_array_elements(quotas) as q;
```


### Get the available zone info of each region

```sql
select
  name,
  zone_names
from
  gcp_compute_region;
```


### Count the available zone in each region

```sql
select
  name,
  count(z) as zone_count
from
  gcp_compute_region,
  jsonb_array_elements(zone_names) as z
group by
  name;
```
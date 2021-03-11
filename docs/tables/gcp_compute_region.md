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


### Get the quota info for a region (us-west1)

```sql
select
  name,
  q -> 'metric' as quota_metric,
  q -> 'limit' as quota_limit
from
  gcp_compute_region,
  jsonb_array_elements(quotas) as q
where
  name = 'us-west1'
order by
  quota_metric;
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
  jsonb_array_length(zone_names) as zone_count
from
  gcp_compute_region;
```

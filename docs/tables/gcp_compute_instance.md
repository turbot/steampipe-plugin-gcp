# Table:  gcp_compute_instance

A GCP Compute Instance is a virtual machine (VM) hosted on Google's infrastructure.

## Examples

### Instance count in each availability zone

```sql
select
  zone_name,
  count(*)
from
  gcp_compute_instance
group by
  zone_name
order by
  count desc;
```


### Count the number of instances by instance type
```sql
select
  machine_type_name,
  count(*) as count
from
  gcp_compute_instance
group by
  machine_type_name
order by
  count desc;
```


### List of instances without application label
```sql
select
  name,
  tags
from
  gcp_compute_instance
where
  tags -> 'application' is null;
```




### List instances having deletion protection feature disabled

```sql
select
  name,
  deletion_protection
from
  gcp_compute_instance
where
  not deletion_protection;
```





### List the disk stats attached to the instances

```sql
select
  name,
  count(d) as num_disks,
  sum( (d ->> 'diskSizeGb') :: int ) as total_storage
from
  gcp_compute_instance as i,
  jsonb_array_elements(disks) as d
group by
    name;
```


### Find instances with IP in a given CIDR range 

```sql
select
  name,
  nic ->> 'networkIP' as ip_address
from
  gcp_compute_instance as i,
  jsonb_array_elements(network_interfaces) as nic
where
    (nic ->> 'networkIP') :: inet <<= '10.128.0.0/16' ; 
```


### Find instances that have been stopped for more than 30 days
```sql
select
  name,
  status,
  last_stop_timestamp
from 
    gcp_compute_instance
where 
    status = 'TERMINATED'
    and last_stop_timestamp < current_timestamp - interval '30 days' ;
```
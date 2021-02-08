# Table: gcp_compute_url_map

When you configure an HTTP(S) load balancer or Traffic Director, you create a URL map. This URL map directs traffic to one or more of the following destinations based on rules that you define: Default backend service. Non-default backend servic

### Get the default backend service of each url-map

```sql
select
  name,
  id,
  split_part(default_service, '/', 10)
from
  gcp_compute_url_map;
```


### Path matcher info of each url-map

```sql
select
  name,
  id,
  p ->> 'name' as name,
  r ->> 'paths' as paths,
  split_part(r ->> 'service', '/', 10) as servise
from
  gcp_compute_url_map,
  jsonb_array_elements(path_matchers) as p,
  jsonb_array_elements(p -> 'pathRules') as r;
```


### Host rule info of each url-map

```sql
select
  name,
  id,
  p ->> 'hosts' as hosts,
  p ->> 'pathMatcher' as path_matcher
from
  gcp_compute_url_map,
  jsonb_array_elements(host_rules) as p;
```


### List of all global type url-maps

```sql
select
  name,
  id,
  location_type
from
  gcp_compute_url_map
where
  location_type = 'GLOBAL';
```

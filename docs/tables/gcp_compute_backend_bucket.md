# Table: gcp_compute_backend_bucket

Backend buckets allow you to use Google Cloud Storage buckets with HTTP(S) Load Balancing.

## Examples

### Basic info

```sql
select
  name,
  id,
  description,
  bucket_name
from
  gcp_compute_backend_bucket;
```

### List of backend buckets where cloud CDN is not enabled

```sql
select
  name,
  id,
  enable_cdn
from
  gcp_compute_backend_bucket
where
  not enable_cdn;
```

### Backend bucket count per storage bucket

```sql
select
  bucket_name,
  count(*) as backend_bucket_count
from
  gcp_compute_backend_bucket
group by
  bucket_name;
```

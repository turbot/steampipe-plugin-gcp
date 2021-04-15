# Table: gcp_compute_ssl_policy

Secure Sockets Layer (SSL) policies determine what port Transport Layer Security (TLS)
features clients are permitted to use when connecting to load balancers.

## Examples

### Basic info

```sql
select
  name,
  id,
  self_link,
  min_tls_version
from
  gcp_compute_ssl_policy;
```

### List SSL policies with minimum TLS version 1.2 and the MODERN profile

```sql
select
  name,
  id,
  min_tls_version
from
  gcp_compute_ssl_policy
where
  min_tls_version = 'TLS_1_2'
  and profile = 'MODERN';
```

### List SSL policies with the RESTRICTED profile

```sql
select
  name,
  id,
  profile
from
  gcp_compute_ssl_policy
where
  profile = 'RESTRICTED';
```

### List SSL policies with weak cipher suites

```sql
select
  name,
  id,
  enabled_feature
from
  gcp_compute_ssl_policy,
  jsonb_array_elements_text(enabled_features) as enabled_feature
where
  profile = 'CUSTOM'
  and enabled_feature in('TLS_RSA_WITH_AES_128_GCM_SHA256', 'TLS_RSA_WITH_AES_256_GCM_SHA384', 'TLS_RSA_WITH_AES_128_CBC_SHA', 'TLS_RSA_WITH_AES_256_CBC_SHA', 'TLS_RSA_WITH_3DES_EDE_CBC_SHA');
```

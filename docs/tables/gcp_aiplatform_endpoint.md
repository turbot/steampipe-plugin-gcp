# Table: gcp_aiplatform_endpoint

AI Platform is a managed service that enables you to easily build machine learning models, that work on any type of data, of any size.

### Basic info

```sql
select
  name,
  display_name,
  create_time,
  network
from
  gcp_aiplatform_endpoint;
```


### List endpoints that are exposed via private service connect

```sql
select
  name,
  display_name,
  create_time,
  enable_private_service_connect
from
  gcp_aiplatform_endpoint
where
  enable_private_service_connect;
```

### List endpoints created in the last 30 days

```sql
select
  name,
  display_name,
  network,
  create_time,
  update_time
from
  gcp_aiplatform_endpoint
where
  create_time >= now() - interval '30' day;
```

### Get customer managed key details of endpoings

```sql
select
  name,
  create_time,
  encryption_spec ->> 'KmsKeyName' as kms_key_name
from
  gcp_aiplatform_endpoint;
```

### Get prediction request response config of endpoints

```sql
select
  name,
  network,
  predict_request_response_logging_config ->> 'Enabled' as enabled,
  predict_request_response_logging_config ->> 'SamplingRate' as sampling_rate,
  predict_request_response_logging_config ->> 'BigqueryDestination' as bigquery_destination
from
  gcp_aiplatform_endpoint;
```
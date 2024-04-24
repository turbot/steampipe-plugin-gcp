
# Test scenarios for the table `gcp_bigquery_dataset` with rate limiter configuration with Single Connection (Over 15,000 Resources)

## List/Get API Call Issues and Observations

### Error Scenario 1 :x:

#### Rate Limiter Settings
- FillRate: 100
- BucketSize: 100

**Query**:
```sql
select * from gcp_large_001.gcp_bigquery_dataset;
```
**Error**:
```
Error: gcp_large_001: table 'gcp_bigquery_dataset' column 'default_partition_expiration_ms' requires hydrate data from getBigQueryDataset, which failed with error googleapi: Error 403: Exceeded rate limits: too many API requests per user per method for this user_method. For more information, see https://cloud.google.com/bigquery/docs/troubleshoot-quotas, rateLimitExceeded.
```
**Stats**:
- Time: 22.9s
- Rows fetched: 6,000
- Hydrate calls: 6,000

---

### Error Scenario 2 :x:

#### Rate Limiter Settings
- FillRate: 100
- BucketSize: 100
- MaxConcurrency: 200
- Scope: `["connection", "service", "action"]`
- Where: "service = 'bigquery' and action = 'datasets.list'"

**Query**:
```sql
select * from gcp_large_001.gcp_bigquery_dataset
```
**Error**:
```
Error: gcp_large_001: table 'gcp_bigquery_dataset' column 'description' requires hydrate data from getBigQueryDataset, which failed with error googleapi: Error 403: Exceeded rate limits: too many API requests per user per method for this user_method. For more information, see https://cloud.google.com/bigquery/docs/troubleshoot-quotas, rateLimitExceeded.
```
**Stats**:
- Time: 17.4s
- Rows fetched: 6,000
- Hydrate calls: 6,000

---

### Error Scenario 3 :x:

#### Rate Limiter Settings
- FillRate: 90
- BucketSize: 90
- MaxConcurrency: 90
- Scope: `["connection", "service", "action"]`
- Where: "service = 'bigquery' and (action = 'datasets.list' or action = 'datasets.get')"

**Query**:
```sql
select * from gcp_large_001.gcp_bigquery_dataset
```
**Error**:
```
Error: gcp_large_001: googleapi: Error 403: Exceeded rate limits: too many API requests per user per method for this user_method. For more information, see https://cloud.google.com/bigquery/docs/troubleshoot-quotas, rateLimitExceeded.
```
**Stats**:
- Time: 184.2s
- Rows fetched: 14,000
- Hydrate calls: 13,553

---

### Error-Free Scenario :white_check_mark:

#### Rate Limiter Settings
- Name: "gcp_bigquery_list_datasets"
- FillRate: 100
- BucketSize: 100
- MaxConcurrency: 50
- Scope: `["connection", "service", "action"]`
- Where: "service = 'bigquery' and (action = 'datasets.list' or action = 'datasets.get')"

**Query**:
```sql
select * from gcp_large_001.gcp_bigquery_dataset
```
**Stats**:
- Time: 201.8s
- Rows fetched: 15,000
- Hydrate calls: 15,000
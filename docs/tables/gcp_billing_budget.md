# Table: gcp_billing_budget

Cloud Billing budgets monitors all your Google Cloud charges in one place. A budget enables you to track your actual Google Cloud costs against your planned costs. After you've set a budget amount, you set budget alert threshold rules that are used to trigger email notifications

**_Please note_**: This table requires the `billing.viewer` permission to retrieve billing budget details.

## Examples

### Basic info

```sql
select
  name,
  billing_account
  display_name,
  specified_amount ->> 'units' as units,
  specified_amount ->> 'currencyCode' as currency_code,
  project,
  location
from
  gcp_billing_budget;
```

### Get threshold rules to trigger alerts for each budget

```sql
select
  name,
  display_name,
  ((threshold_rule ->> 'thresholdPercent')::numeric) * 100 || '%' as threshold_percent,
  threshold_rule ->> 'spendBasis' as spend_basis
from
  gcp_billing_budget,
  jsonb_array_elements(threshold_rules) as threshold_rule;
```

### Get filters limiting the scope of the cost to calculate budget

```sql
select
  name,
  display_name,
  string_agg(p, ', ') as applies_to_projects,
  specified_amount ->> 'units' as units,
  specified_amount ->> 'currencyCode' as currency_code,
  budget_filter ->> 'calendarPeriod' as budget_calendar_period,
  budget_filter ->> 'creditTypesTreatment' as budget_credit_types_treatment
from
  gcp_billing_budget,
  jsonb_array_elements_text(budget_filter -> 'projects') as p
group by
  name,
  display_name,
  budget_filter,
  specified_amount;
```
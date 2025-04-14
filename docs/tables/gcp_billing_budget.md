---
title: "Steampipe Table: gcp_billing_budget - Query GCP Billing Budgets using SQL"
description: "Allows users to query Billing Budgets in Google Cloud Platform, specifically the budget amount, associated projects, and alert thresholds, providing insights into budget management and cost control."
folder: "Billing"
---

# Table: gcp_billing_budget - Query GCP Billing Budgets using SQL

A Billing Budget in Google Cloud Platform is a tool that allows you to set a custom cost threshold for your Google Cloud usage. It can be associated with one or more projects and can be configured to send alerts when the usage approaches or exceeds the set budget. This helps in managing costs, avoiding overspending, and keeping track of your cloud resource consumption.

## Table Usage Guide

The `gcp_billing_budget` table provides insights into Billing Budgets within Google Cloud Platform. As a finance or operations manager, explore budget-specific details through this table, including budget amounts, associated projects, and alert thresholds. Utilize it to manage costs, monitor spending, and ensure that resource usage is within the set budgets.

**Important Notes**
- This table requires the `billing.viewer` permission to retrieve billing account details.

## Examples

### Basic info
Explore your Google Cloud Platform's budget details to gain insights into the specified amounts, including units and currency codes, for each project and location. This can help you manage your resources more effectively and keep track of your spending.

```sql+postgres
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

```sql+sqlite
select
  name,
  billing_account,
  display_name,
  json_extract(specified_amount, '$.units') as units,
  json_extract(specified_amount, '$.currencyCode') as currency_code,
  project,
  location
from
  gcp_billing_budget;
```

### Get threshold rules to trigger alerts for each budget
Explore the budget alert rules to understand when each budget will trigger an alert based on a certain spending threshold. This is useful to manage and control spending within your budget limits.

```sql+postgres
select
  name,
  display_name,
  ((threshold_rule ->> 'thresholdPercent')::numeric) * 100 || '%' as threshold_percent,
  threshold_rule ->> 'spendBasis' as spend_basis
from
  gcp_billing_budget,
  jsonb_array_elements(threshold_rules) as threshold_rule;
```

```sql+sqlite
select
  name,
  display_name,
  (json_extract(threshold_rule.value, '$.thresholdPercent') * 100) || '%' as threshold_percent,
  json_extract(threshold_rule.value, '$.spendBasis') as spend_basis
from
  gcp_billing_budget,
  json_each(threshold_rules) as threshold_rule;
```

### Get filters limiting the scope of the cost to calculate budget
This query is useful for gaining insights into your Google Cloud Platform (GCP) billing budget. It allows you to identify the filters that limit the cost scope for budget calculations, providing a better understanding of your spending limits and how they are distributed across different projects.

```sql+postgres
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

```sql+sqlite
select
  name,
  display_name,
  group_concat(p.value, ', ') as applies_to_projects,
  json_extract(specified_amount, '$.units') as units,
  json_extract(specified_amount, '$.currencyCode') as currency_code,
  json_extract(budget_filter, '$.calendarPeriod') as budget_calendar_period,
  json_extract(budget_filter, '$.creditTypesTreatment') as budget_credit_types_treatment
from
  gcp_billing_budget,
  json_each(json_extract(budget_filter, '$.projects')) as p
group by
  name,
  display_name,
  budget_filter,
  specified_amount;
```
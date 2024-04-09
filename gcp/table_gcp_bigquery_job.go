package gcp

import (
	"context"
	"strings"

	"github.com/turbot/go-kit/types"
	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	"google.golang.org/api/bigquery/v2"
)

//// TABLE DEFINITION

func tableGcpBigQueryJob(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_bigquery_job",
		Description: "GCP BigQuery Job",
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("job_id"),
			Hydrate:    getBigQueryJob,
		},
		List: &plugin.ListConfig{
			Hydrate: listBigQueryJobs,
			Tags:    map[string]string{"service": "bigquery", "action": "jobs.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "Unique opaque ID of the job.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "job_id",
				Description: "The ID of the job.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("JobReference.JobId"),
			},
			{
				Name:        "state",
				Description: "Running state of the job. When the state is DONE, errorResult can be checked to determine whether the job succeeded or failed.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Status.State"),
			},
			{
				Name:        "self_link",
				Description: "An URL that can be used to access the resource again.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBigQueryJob,
			},
			{
				Name:        "completion_ratio",
				Description: "Job progress (0.0 -> 1.0) for LOAD and EXTRACT jobs.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Statistics.CompletionRatio"),
			},
			{
				Name:        "creation_time",
				Description: "Creation time of this job.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Statistics.CreationTime").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "end_time",
				Description: "End time of this job.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Statistics.EndTime").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "etag",
				Description: "A hash of the resource.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBigQueryJob,
			},
			{
				Name:        "kind",
				Description: "The type of the resource.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "num_child_jobs",
				Description: "Number of child jobs executed.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Statistics.NumChildJobs"),
			},
			{
				Name:        "parent_job_id",
				Description: "If this is a child job, the id of the parent.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Statistics.ParentJobId"),
			},
			{
				Name:        "reservation_id",
				Description: "Name of the primary reservation assigned to this job.",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Statistics.ReservationId"),
			},
			{
				Name:        "start_time",
				Description: "Start time of this job.",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Statistics.StartTime").Transform(transform.UnixMsToTimestamp),
			},
			{
				Name:        "total_bytes_processed",
				Description: "Use the bytes processed in the query statistics instead.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Statistics.TotalBytesProcessed"),
			},
			{
				Name:        "total_slot_ms",
				Description: "Slot-milliseconds for the job.",
				Type:        proto.ColumnType_INT,
				Transform:   transform.FromField("Statistics.TotalSlotMs"),
			},
			{
				Name:        "user_email",
				Description: "Email address of the user who ran the job.",
				Type:        proto.ColumnType_STRING,
				Hydrate:     getBigQueryJob,
			},
			{
				Name:        "configuration",
				Description: "Describes the job configuration.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBigQueryJob,
			},
			{
				Name:        "error_result",
				Description: "A result object that will be present only if the job has failed.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBigQueryJob,
				Transform:   transform.FromField("Status.ErrorResult"),
			},
			{
				Name:        "errors",
				Description: "he first errors encountered during the running of the job. The final message includes the number of errors that caused the process to stop. Errors here do not necessarily mean that the job has completed or was unsuccessful.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBigQueryJob,
				Transform:   transform.FromField("Status.Errors"),
			},
			{
				Name:        "extract",
				Description: "Statistics for an extract job.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Statistics.Extract"),
			},
			{
				Name:        "labels",
				Description: "The labels associated with this job.",
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBigQueryJob,
				Transform:   transform.FromField("Configuration.Labels"),
			},
			{
				Name:        "load",
				Description: "Statistics for a load job.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Statistics.Load"),
			},
			{
				Name:        "query",
				Description: "Statistics for a query job.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Statistics.Query"),
			},
			{
				Name:        "quota_deferments",
				Description: "Quotas which delayed this job's start time.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Statistics.QuotaDeferments"),
			},
			{
				Name:        "reservation_usage",
				Description: "Job resource usage breakdown by reservation.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Statistics.ReservationUsage"),
			},
			{
				Name:        "row_level_security_statistics",
				Description: "Statistics for row-level security. Present only for query and extract jobs.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Statistics.RowLevelSecurityStatistics"),
			},
			{
				Name:        "script_statistics",
				Description: "Statistics for a child job of a script.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Statistics.ScriptStatistics"),
			},
			{
				Name:        "session_info_template",
				Description: "Information of the session if this job is part of one.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Statistics.SessionInfoTemplate"),
			},
			{
				Name:        "transaction_info",
				Description: "Information of the multi-statement transaction if this job is part of one.",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Statistics.TransactionInfo"),
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("JobReference.JobId"),
			},
			{
				Name:        "tags",
				Description: ColumnDescriptionTags,
				Type:        proto.ColumnType_JSON,
				Hydrate:     getBigQueryJob,
				Transform:   transform.FromField("Configuration.Labels"),
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.From(bigQueryJobAka),
			},

			// GCP standard columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("JobReference.Location"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("JobReference.ProjectId"),
			},
		},
	}
}

//// LIST FUNCTION

func listBigQueryJobs(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	plugin.Logger(ctx).Trace("listBigQueryJobs")

	// Create Service Connection
	service, err := BigQueryService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Max limit isn't mentioned in the documentation
	// Default limit is set as 1000
	pageSize := types.Int64(1000)
	limit := d.QueryContext.Limit
	if d.QueryContext.Limit != nil {
		if *limit < *pageSize {
			pageSize = limit
		}
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	resp := service.Jobs.List(project).MaxResults(*pageSize).AllUsers(true)
	if err := resp.Pages(ctx, func(page *bigquery.JobList) error {
		// apply rate limiting
		d.WaitForListRateLimit(ctx)

		for _, job := range page.Jobs {
			d.StreamListItem(ctx, job)

			// Check if context has been cancelled or if the limit has been hit (if specified)
			// if there is a limit, it will return the number of rows required to reach this limit
			if d.RowsRemaining(ctx) == 0 {
				page.NextPageToken = ""
				return nil
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return nil, nil
}

//// HYDRATE FUNCTIONS

func getBigQueryJob(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Create Service Connection
	service, err := BigQueryService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	var id string
	if h.Item != nil {
		data := jobID(h.Item)
		id = strings.Split(data, ".")[1]
	} else {
		id = d.EqualsQuals["job_id"].GetStringValue()
	}

	// handle empty id in get call
	if id == "" {
		return nil, nil
	}

	resp, err := service.Jobs.Get(project, id).Do()
	if err != nil {
		return nil, err
	}

	return resp, nil
}

//// TRANSFORM FUNCTIONS

func bigQueryJobAka(ctx context.Context, h *transform.TransformData) (interface{}, error) {
	data := jobID(h.HydrateItem)
	projectID := strings.Split(data, ":")[0]
	id := strings.Split(data, ".")[1]

	akas := []string{"gcp://bigquery.googleapis.com/projects/" + projectID + "/jobs/" + id}

	return akas, nil
}

func jobID(item interface{}) string {
	switch item := item.(type) {
	case *bigquery.JobListJobs:
		return item.Id
	case *bigquery.Job:
		return item.Id
	}
	return ""
}

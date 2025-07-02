package gcp

import (
	"context"
	"time"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
	adminreports "google.golang.org/api/admin/reports/v1"
)

//// TABLE DEFINITION

func tableGcpAdminReportsLoginActivity(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_admin_reports_login_activity",
		Description: "GCP Admin Reports API - activité de connexion (login)",
		List: &plugin.ListConfig{
			Hydrate: listGcpAdminReportsLoginActivities,
			KeyColumns: plugin.KeyColumnSlice{
				{Name: "time", Require: plugin.Optional, Operators: []string{">", ">=", "<", "<=", "="}},
				{Name: "actor_email", Require: plugin.Optional},
				{Name: "ip_address", Require: plugin.Optional},
				{Name: "event_name", Require: plugin.Optional},
			},
			Tags: map[string]string{"service": "admin", "product": "reports", "action": "activities.list"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "time",
				Description: "Timestamp of the activity (Id.Time) in RFC3339 format",
				Type:        proto.ColumnType_TIMESTAMP,
				Transform:   transform.FromField("Id.Time"),
			},
			{
				Name:        "actor_email",
				Description: "Email address of the actor (Actor.Email)",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Actor.Email"),
			},
			{
    			Name:        "event_name",
    			Description: "List of event names for this activity",
    			Type:        proto.ColumnType_STRING,
    			Transform:   transform.FromField("Events").Transform(extractFirstEventName),
			},
			{
				Name:        "unique_qualifier",
				Description: "Unique qualifier ID for this activity (Id.UniqueQualifier)",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id.UniqueQualifier"),
			},
			{
				Name:        "application_name",
				Description: "Name of the report application (Id.ApplicationName)",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id.ApplicationName"),
			},
			{
				Name:        "actor_profile_id",
				Description: "Profile ID of the actor (Actor.ProfileId)",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Actor.ProfileId"),
			},
			{
				Name:        "actor_caller_type",
				Description: "Caller type of the actor (Actor.CallerType)",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Actor.CallerType"),
			},
			{
				Name:        "ip_address",
				Description: "IP address associated with the activity (IpAddress)",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("IpAddress"),
			},
			{
				Name:        "events",
				Description: "Full JSON array of detailed events (Events)",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Events"),
			},
			{
				Name:        "title",
				Description: "Concatenation of time and actor email",
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Id.Time").Transform(convertTimeToString).Transform(formatTitleWithActorEmail),
			},
			{
				Name:        "tags",
				Description: "Tags (List of events)",
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Events").Transform(extractEventNames),
			},
		},
	}
}


//// LIST FUNCTION

func listGcpAdminReportsLoginActivities(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
    // Création du service Reports API
    service, err := ReportsService(ctx, d)
    if err != nil {
        plugin.Logger(ctx).Error("gcp_admin_reports_login_activity.list", "service_error", err)
        return nil, err
    }

    call := service.Activities.List("all", "login")

    // If the user supplied a time qualifier, translate it to StartTime/EndTime parameters
	if quals := d.Quals["time"]; quals != nil {
		var startTime, endTime time.Time
		for _, q := range quals.Quals {
			if ts := q.Value.GetTimestampValue(); ts != nil {
				t := ts.AsTime()
				switch q.Operator {
				case "=":
					startTime, endTime = t, t
				case ">":
					startTime = t.Add(time.Nanosecond)
				case ">=":
					startTime = t
				case "<":
					endTime = t
				case "<=":
					endTime = t
				}
			}
		}
		if !startTime.IsZero() {
			call.StartTime(startTime.Format(time.RFC3339))
		}
		if !endTime.IsZero() {
			call.EndTime(endTime.Format(time.RFC3339))
		}
	}

	// Pagination setup
	pageToken := ""
	const apiMaxPageSize = 1000

	// Determine initial page size based on SQL LIMIT
	var initialPageSize int64 = apiMaxPageSize
	if limit := d.QueryContext.Limit; limit != nil && *limit < initialPageSize {
		initialPageSize = *limit
	}
	call.MaxResults(initialPageSize)

	for {
		if pageToken != "" {
			call.PageToken(pageToken)
		}
		resp, err := call.Do()
		if err != nil {
			plugin.Logger(ctx).Error("gcp_admin_reports_login_activity.list", "api_error", err)
			return nil, err
		}
		// Stream items
		for _, activity := range resp.Items {
			d.StreamListItem(ctx, activity)
			if d.RowsRemaining(ctx) == 0 {
				return nil, nil
			}
		}
		// Check for next page
		if resp.NextPageToken == "" {
			break
		}
		pageToken = resp.NextPageToken

		// Adjust next page size for remaining rows
		if limit := d.QueryContext.Limit; limit != nil {
			remaining := d.RowsRemaining(ctx)
			if remaining > 0 && remaining < apiMaxPageSize {
				call.MaxResults(int64(remaining))
			} else {
				call.MaxResults(apiMaxPageSize)
			}
		} else {
			call.MaxResults(apiMaxPageSize)
		}
	}

	return nil, nil
}


//// TRANSFORM FUNCTIONS 
func extractEventNames(_ context.Context, d *transform.TransformData) (interface{}, error) {
	activity, ok := d.HydrateItem.(*adminreports.Activity)
	if !ok {
		return nil, nil
	}
	if activity.Events == nil {
		return nil, nil
	}
	names := []string{}
	for _, e := range activity.Events {
		if e.Name != "" {
			names = append(names, e.Name)
		}
	}
	return names, nil
}

func extractFirstEventName(_ context.Context, d *transform.TransformData) (interface{}, error) {
    events, ok := d.Value.([]*adminreports.ActivityEvents)
    if !ok || len(events) == 0 {
        return "", nil
    }
    return events[0].Name, nil
}

func convertTimeToString(_ context.Context, d *transform.TransformData) (interface{}, error) {
	activity, ok := d.HydrateItem.(*adminreports.Activity)
	if !ok {
		return "", nil
	}
	if activity.Id == nil || activity.Id.Time == "" {
		return "", nil
	}
	return activity.Id.Time, nil
}

func formatTitleWithActorEmail(_ context.Context, d *transform.TransformData) (interface{}, error) {
	timeStr, ok := d.Value.(string)
	if !ok {
		return nil, nil
	}
	activity, ok := d.HydrateItem.(*adminreports.Activity)
	if !ok {
		return timeStr, nil
	}
	if activity.Actor == nil || activity.Actor.Email == "" {
		return timeStr, nil
	}
	return timeStr + " - " + activity.Actor.Email, nil
}
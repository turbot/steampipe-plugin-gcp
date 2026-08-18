package gcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/turbot/steampipe-plugin-sdk/v6/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v6/plugin/transform"
	"google.golang.org/api/iterator"

	recommender "cloud.google.com/go/recommender/apiv1"
	recommenderpb "cloud.google.com/go/recommender/apiv1/recommenderpb"
)

//// TABLE DEFINITION

// Reference: https://docs.cloud.google.com/policy-intelligence/docs/service-account-insights#get-a-single-service-account-insight
func tableGcpRecommenderInsights(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_recommender_insight",
		Description: "GCP Recommender Insights",
		List: &plugin.ListConfig{
			Hydrate: listGcpRecommenderInsights,
			Tags:    map[string]string{"service": "recommender", "action": "recommender.listInsight"},
		},
		Get: &plugin.GetConfig{
			KeyColumns: plugin.SingleColumn("project"),
			Hydrate:    getGcpRecommenderInsights,
			Tags:       map[string]string{"service": "recommender", "action": "recommender.listInsight"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "category",
				Description: "The category for IAM insights is always SECURITY.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "content",
				Description: "Reports the last time the service account was authenticated.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "description",
				Description: "A human-readable summary of the insight.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "etag",
				Description: "A unique identifier for the current state of an insight. Each time the insight changes, a new etag value is assigned.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_refresh_time",
				Description: "The date when the insight was last refreshed, which indicates the freshness of the data used to generate the insight.",
				Type:        proto.ColumnType_TIMESTAMP,
			},
			{
				Name:        "name",
				Description: "The name of the insight.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "observation_period",
				Description: "The time period leading up to the insight. The source data used to generate the insight ends at lastRefreshTime and begins at lastRefreshTime minus observationPeriod..",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "state_info",
				Description: "Insights go through multiple state transitions after they are proposed.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "target_resources",
				Description: "The full resource name of the project that the insight is for.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_account_email",
				Description: "The email address of the service account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_account_id",
				Description: "The unique numeric ID of the service account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "last_authenticated_time",
				Description: "The most recent time that the service account was authenticated. If the service account does not have any recorded authentications, this field is not included..",
				Type:        proto.ColumnType_TIMESTAMP,
			},

			// standard steampipe columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getRecommenderInsightsTurbotData,
			},
			{
				Name:        "akas",
				Description: ColumnDescriptionAkas,
				Type:        proto.ColumnType_JSON,
				Transform:   transform.FromField("Name"),
			},

			// standard gcp columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromConstant("global"),
			},
			{
				Name:        "project",
				Description: ColumnDescriptionProject,
				Type:        proto.ColumnType_STRING,
				Hydrate:     getProject,
				Transform:   transform.FromValue(),
			},
		},
	}
}

//// FETCH FUNCTIONS

func listGcpRecommenderInsights(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := recommender.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Get project details
	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)
	plugin.Logger(ctx).Trace("listGcpRecommenderInsights", "GCP_PROJECT: ", project)

	parent := fmt.Sprintf("projects/%s/locations/global/insightTypes/google.iam.serviceAccount.Insight", projectId)

	req := &recommenderpb.ListInsightsRequest{
		// See https://pkg.go.dev/cloud.google.com/go/recommender/apiv1/recommenderpb#ListInsightsRequest
		Parent: parent,
	}
	it := client.ListInsights(ctx, req)

	// apply rate limiting
	d.WaitForListRateLimit(ctx)

	for {
		resp, err := it.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		d.StreamListItem(ctx, resp)
	}

	return nil, nil
}

func getGcpRecommenderInsights(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	client, err := recommender.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	// Get project id
	projectId := d.EqualsQuals["project"].GetStringValue()
	plugin.Logger(ctx).Trace("getGcpRecommenderInsights", "GCP_PROJECT: ", projectId)

	parent := fmt.Sprintf("projects/%s/locations/global/insightTypes/google.iam.serviceAccount.Insight", projectId)

	req := &recommenderpb.ListInsightsRequest{
		// See https://pkg.go.dev/cloud.google.com/go/recommender/apiv1/recommenderpb#ListInsightsRequest
		Parent: parent,
	}
	it := client.ListInsights(ctx, req)

	// apply rate limiting
	d.WaitForListRateLimit(ctx)

	for {
		resp, err := it.Next()

		if err == iterator.Done {
			break
		}

		if err != nil {
			return nil, err
		}

		d.StreamListItem(ctx, resp)
	}

	return nil, nil
}

func getRecommenderInsightsTurbotData(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	// Get the resource title
	title := strings.ToUpper(project) + " Recommender Insights"

	// Build resource aka
	akas := []string{"gcp://cloudresourcemanager.googleapis.com/projects/" + project + "/iamPolicy"}

	// Mapping all turbot defined properties
	turbotData := map[string]interface{}{
		"Akas":  akas,
		"Title": title,
	}

	return turbotData, nil
}

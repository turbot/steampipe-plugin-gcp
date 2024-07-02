package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/v5/grpc/proto"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin/transform"
)

//// TABLE DEFINITION

// We can have only one Application per project. The App ID would be the Project ID.
// https://cloud.google.com/appengine/docs/flexible/managing-projects-apps-billing#:~:text=Important%3A%20Each%20Google%20Cloud%20project,of%20your%20App%20Engine%20application
func tableGcpAppEngineApplication(ctx context.Context) *plugin.Table {
	return &plugin.Table{
		Name:        "gcp_app_engine_application",
		Description: "GCP App Engine Application",
		List: &plugin.ListConfig{
			Hydrate: getAppEngineApplication,
			Tags:    map[string]string{"service": "appengine", "action": "applications.get"},
		},
		Columns: []*plugin.Column{
			{
				Name:        "id",
				Description: "Identifier of the Application resource. This identifier is equivalent to the project ID of the Google Cloud Platform project where you want to deploy your application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "name",
				Description: "Full path to the Application resource in the API.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "service_account",
				Description: "The service account associated with the application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "auth_domain",
				Description: "Google Apps authentication domain that controls which users can access this application.Defaults to open access for any Google Account.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "code_bucket",
				Description: "Google Cloud Storage bucket that can be used for storing files associated with this application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "database_type",
				Description: "The type of the Cloud Firestore or Cloud Datastore database associated with this application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_bucket",
				Description: "Google Cloud Storage bucket that can be used by this application to store content.@OutputOnly.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_cookie_expiration",
				Description: "Cookie expiration policy for this application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "default_hostname",
				Description: "Hostname used to reach this application, as resolved by App Engine.@OutputOnly.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "gcr_domain",
				Description: "The Google Container Registry domain used for storing managed build docker images for this application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "serving_status",
				Description: "Serving status of this application.",
				Type:        proto.ColumnType_STRING,
			},
			{
				Name:        "dispatch_rules",
				Description: "HTTP path dispatch rules for requests to the application that do not explicitly target a service or version.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "feature_settings",
				Description: "The feature specific settings to be used in the application.",
				Type:        proto.ColumnType_JSON,
			},
			{
				Name:        "iap",
				Description: "Identity-Aware Proxy.",
				Type:        proto.ColumnType_JSON,
			},

			// Steampipe standard columns
			{
				Name:        "title",
				Description: ColumnDescriptionTitle,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("Name"),
			},

			// GCP standard columns
			{
				Name:        "location",
				Description: ColumnDescriptionLocation,
				Type:        proto.ColumnType_STRING,
				Transform:   transform.FromField("LocationId"),
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

//// LIST FUNCTION

func getAppEngineApplication(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {

	// Create Service Connection
	service, err := AppEngineService(ctx, d)
	if err != nil {
		return nil, err
	}

	// Get project details

	projectId, err := getProject(ctx, d, h)
	if err != nil {
		return nil, err
	}
	project := projectId.(string)

	// In Google Cloud Platform (GCP), the structure is such that each project
	// can contain only one App Engine application. This means the number of
	// App Engine applications is directly related to the number of projects you have.

	// Each GCP project can support various services and features from Google Cloud,
	// but when it comes to App Engine, it is restricted to one application per project.

	// If multiple applications are needed, you will need to create additional GCP projects,
	// with one project for each App Engine application you wish to deploy. This approach
	// Available APIs: https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1/apps

	resp, err := service.Apps.Get(project).Do()
	if err != nil {
		return nil, err
	}

	d.StreamListItem(ctx, resp)

	return resp, nil
}

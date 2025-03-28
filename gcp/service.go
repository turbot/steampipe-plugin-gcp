package gcp

import (
	"context"

	aiplatform "cloud.google.com/go/aiplatform/apiv1"
	redis "cloud.google.com/go/redis/apiv1"
	rediscluster "cloud.google.com/go/redis/cluster/apiv1"
	"github.com/turbot/steampipe-plugin-sdk/v5/plugin"
	"google.golang.org/api/accessapproval/v1"
	"google.golang.org/api/alloydb/v1"
	"google.golang.org/api/apikeys/v2"
	"google.golang.org/api/appengine/v1"
	"google.golang.org/api/artifactregistry/v1"
	"google.golang.org/api/bigquery/v2"
	"google.golang.org/api/bigtableadmin/v2"
	"google.golang.org/api/billingbudgets/v1"
	"google.golang.org/api/cloudasset/v1"
	"google.golang.org/api/cloudbilling/v1"
	"google.golang.org/api/cloudfunctions/v2"
	"google.golang.org/api/cloudidentity/v1"
	"google.golang.org/api/cloudkms/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/composer/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/container/v1"
	"google.golang.org/api/dataplex/v1"
	"google.golang.org/api/dataproc/v1"
	"google.golang.org/api/dns/v1"
	"google.golang.org/api/essentialcontacts/v1"
	"google.golang.org/api/firestore/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/logging/v2"
	"google.golang.org/api/metastore/v1"
	"google.golang.org/api/monitoring/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/pubsub/v1"
	run1 "google.golang.org/api/run/v1"
	"google.golang.org/api/run/v2"
	"google.golang.org/api/secretmanager/v1"
	"google.golang.org/api/serviceusage/v1"
	"google.golang.org/api/storage/v1"
	"google.golang.org/api/vpcaccess/v1"

	computeBeta "google.golang.org/api/compute/v0.beta"
	sqladmin "google.golang.org/api/sqladmin/v1beta4"
)

// AccessApprovalService returns the service connection for GCP Project AccessApproval service
func AccessApprovalService(ctx context.Context, d *plugin.QueryData) (*accessapproval.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "AccessApprovalService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*accessapproval.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := accessapproval.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

type AIplatfromServiceClients struct {
	Endpoint *aiplatform.EndpointClient
	Dataset  *aiplatform.DatasetClient
	Index    *aiplatform.IndexClient
	Job      *aiplatform.JobClient
	Model    *aiplatform.ModelClient
	Notebook *aiplatform.NotebookClient
}

func AIService(ctx context.Context, d *plugin.QueryData, clientType string) (*AIplatfromServiceClients, error) {
	// have we already created and cached the service?
	matrixLocation := d.EqualsQualString(matrixKeyLocation)

	// Default to us-central1 for building the supported locations for the resources like Endpoint, Dataset, Index, Job etc...
	if matrixLocation == "" {
		matrixLocation = "us-central1"
	}

	serviceCacheKey := "AIService" + matrixLocation + clientType
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*AIplatfromServiceClients), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)
	opts = append(opts, option.WithEndpoint(matrixLocation+"-aiplatform.googleapis.com:443"))

	clients := &AIplatfromServiceClients{}

	switch clientType {
	case "Endpoint":
		svc, err := aiplatform.NewEndpointClient(ctx, opts...)
		if err != nil {
			return nil, err
		}
		clients.Endpoint = svc
		return clients, nil
	case "Dataset":
		svc, err := aiplatform.NewDatasetClient(ctx, opts...)
		if err != nil {
			return nil, err
		}
		clients.Dataset = svc
		return clients, nil
	case "Index":
		svc, err := aiplatform.NewIndexClient(ctx, opts...)
		if err != nil {
			return nil, err
		}
		clients.Index = svc
		return clients, nil
	case "Job":
		svc, err := aiplatform.NewJobClient(ctx, opts...)
		if err != nil {
			return nil, err
		}
		clients.Job = svc
		return clients, nil
	case "Model":
		svc, err := aiplatform.NewModelClient(ctx, opts...)
		if err != nil {
			return nil, err
		}
		clients.Model = svc
		return clients, nil
	case "Notebook":
		svc, err := aiplatform.NewNotebookClient(ctx, opts...)
		if err != nil {
			return nil, err
		}
		clients.Notebook = svc
		return clients, nil
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, clients)
	return clients, nil
}

// AlloyDBService returns the service connection for GCP Alloy DB service
func AlloyDBService(ctx context.Context, d *plugin.QueryData) (*alloydb.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "AlloyDBservice"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*alloydb.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := alloydb.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

func APIKeysService(ctx context.Context, d *plugin.QueryData) (*apikeys.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "APIKeysService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*apikeys.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := apikeys.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// AppEngineService returns the service connection for GCP App Engine service
func AppEngineService(ctx context.Context, d *plugin.QueryData) (*appengine.APIService, error) {
	// have we already created and cached the service?
	serviceCacheKey := "BillingBudgetsService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*appengine.APIService), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := appengine.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// BillingBudgetsService returns the service connection for GCP Billing Budgets service
func BillingBudgetsService(ctx context.Context, d *plugin.QueryData) (*billingbudgets.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "BillingBudgetsService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*billingbudgets.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := billingbudgets.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// BillingService returns the service connection for GCP Billing service
func BillingService(ctx context.Context, d *plugin.QueryData) (*cloudbilling.APIService, error) {
	// have we already created and cached the service?
	serviceCacheKey := "BillingService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*cloudbilling.APIService), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := cloudbilling.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// BigQueryService returns the service connection for GCP BigQueryService service
func BigQueryService(ctx context.Context, d *plugin.QueryData) (*bigquery.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "BigQueryService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*bigquery.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := bigquery.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// ArtifactRegistryService returns the service connection for GCP ArtifactRegistry service
func ArtifactRegistryService(ctx context.Context, d *plugin.QueryData) (*artifactregistry.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "ArtifactRegistryService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*artifactregistry.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := artifactregistry.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// BigtableAdminService returns the service connection for GCP Bigtable Admin service
func BigtableAdminService(ctx context.Context, d *plugin.QueryData) (*bigtableadmin.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "BigtableAdminService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*bigtableadmin.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := bigtableadmin.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// CloudResourceManagerService returns the service connection for GCP Cloud Resource Manager service
func CloudResourceManagerService(ctx context.Context, d *plugin.QueryData) (*cloudresourcemanager.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "CloudResourceManagerService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*cloudresourcemanager.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := cloudresourcemanager.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// CloudRunService returns the service connection for GCP Cloud Run service
func CloudRunService(ctx context.Context, d *plugin.QueryData) (*run.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "CloudRunService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*run.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := run.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// CloudRunService returns the service connection for GCP Cloud Run service
func CloudRunServiceV1(ctx context.Context, d *plugin.QueryData) (*run1.APIService, error) {
	// have we already created and cached the service?
	serviceCacheKey := "CloudRunServiceV1"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*run1.APIService), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := run1.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// DataplexService returns the service connection for GCP Dataplex service
func DataplexService(ctx context.Context, d *plugin.QueryData) (*dataplex.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "DataplexService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*dataplex.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := dataplex.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// EssentialContactService returns the service connection for GCP Cloud Organization Essential Contacts
func EssentialContactService(ctx context.Context, d *plugin.QueryData) (*essentialcontacts.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "EssentialContactService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*essentialcontacts.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := essentialcontacts.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// CloudSQLAdminService returns the service connection for GCP Cloud SQL Admin service
func CloudSQLAdminService(ctx context.Context, d *plugin.QueryData) (*sqladmin.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "CloudSQLAdminService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*sqladmin.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := sqladmin.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// ComputeBetaService returns the service connection for GCP Compute service beta version
func ComputeBetaService(ctx context.Context, d *plugin.QueryData) (*computeBeta.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "ComputeBetaService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*computeBeta.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := computeBeta.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// ComputeService returns the service connection for GCP Compute service
func ComputeService(ctx context.Context, d *plugin.QueryData) (*compute.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "ComputeService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*compute.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := compute.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// ComposerService returns the service connection for GCP Composer service
func ComposerService(ctx context.Context, d *plugin.QueryData) (*composer.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "ComposerService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*composer.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := composer.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// DataprocService returns the service connection for GCP Dataproc service
func DataprocService(ctx context.Context, d *plugin.QueryData) (*dataproc.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "DataprocService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*dataproc.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := dataproc.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// DataprocService returns the service connection for GCP Dataproc service
func DataprocMetastoreService(ctx context.Context, d *plugin.QueryData) (*metastore.APIService, error) {
	// have we already created and cached the service?
	serviceCacheKey := "DataprocMetastoreService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*metastore.APIService), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := metastore.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// ContainerService returns the service connection for GCP Container service
func ContainerService(ctx context.Context, d *plugin.QueryData) (*container.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "ContainerService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*container.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := container.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// CloudFunctionsService returns the service connection for GCP Cloud Functions service
func CloudFunctionsService(ctx context.Context, d *plugin.QueryData) (*cloudfunctions.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "CloudFunctionsService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*cloudfunctions.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := cloudfunctions.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// CloudIdentityService returns the service connection for GCP Identity service
func CloudIdentityService(ctx context.Context, d *plugin.QueryData) (*cloudidentity.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "CloudIdentityService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*cloudidentity.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := cloudidentity.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// CloudAssetService returns the service connection for GCP Asset Service
func CloudAssetService(ctx context.Context, d *plugin.QueryData) (*cloudasset.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "CloudIdentityService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*cloudasset.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := cloudasset.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// DnsService returns the service connection for GCP DNS service
func DnsService(ctx context.Context, d *plugin.QueryData) (*dns.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "DnsService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*dns.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := dns.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// FirestoreDatabaseService returns the service connection for GCP Firestore service
func FirestoreDatabaseService(ctx context.Context, d *plugin.QueryData) (*firestore.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "FirestoreDatabaseService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*firestore.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := firestore.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// IAMService returns the service connection for GCP IAM service
func IAMService(ctx context.Context, d *plugin.QueryData) (*iam.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "IAMService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*iam.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := iam.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// LoggingService returns the service connection for GCP Logging service
func LoggingService(ctx context.Context, d *plugin.QueryData) (*logging.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "LoggingService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*logging.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := logging.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// MonitoringService returns the service connection for GCP Monitoring service
func MonitoringService(ctx context.Context, d *plugin.QueryData) (*monitoring.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "MonitoringService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*monitoring.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := monitoring.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// PubsubService returns the service connection for GCP Pub/Sub service
func PubsubService(ctx context.Context, d *plugin.QueryData) (*pubsub.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "PubsubService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*pubsub.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := pubsub.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// ServiceUsageService returns the service connection for GCP Service Usage service
func ServiceUsageService(ctx context.Context, d *plugin.QueryData) (*serviceusage.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "ServiceUsageService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*serviceusage.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := serviceusage.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// StorageService returns the service connection for GCP Storage service
func StorageService(ctx context.Context, d *plugin.QueryData) (*storage.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "StorageService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*storage.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := storage.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// KMSService returns the service connection for GCP KMS service
func KMSService(ctx context.Context, d *plugin.QueryData) (*cloudkms.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "KMSService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*cloudkms.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := cloudkms.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// RedisService returns the service connection for GCP Redis service
func RedisService(ctx context.Context, d *plugin.QueryData) (*redis.CloudRedisClient, error) {
	// have we already created and cached the service?
	serviceCacheKey := "RedisService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*redis.CloudRedisClient), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := redis.NewCloudRedisClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// RedisClusterService returns the service connection for GCP Memorystore for Redis Cluster service
func RedisClusterService(ctx context.Context, d *plugin.QueryData) (*rediscluster.CloudRedisClusterClient, error) {
	// have we already created and cached the service?
	serviceCacheKey := "RedisClusterService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*rediscluster.CloudRedisClusterClient), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := rediscluster.NewCloudRedisClusterClient(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

func SecretManagerService(ctx context.Context, d *plugin.QueryData) (*secretmanager.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "SecretManagerService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*secretmanager.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := secretmanager.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

func VPCAccessService(ctx context.Context, d *plugin.QueryData) (*vpcaccess.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "VPCAccessService"
	if cachedData, ok := d.ConnectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*vpcaccess.Service), nil
	}

	// To get config arguments from plugin config file
	opts := setSessionConfig(ctx, d.Connection)

	// so it was not in cache - create service
	svc, err := vpcaccess.NewService(ctx, opts...)
	if err != nil {
		return nil, err
	}

	d.ConnectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

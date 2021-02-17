package gcp

import (
	"context"

	"github.com/turbot/steampipe-plugin-sdk/connection"
	"google.golang.org/api/cloudfunctions/v1"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/iam/v1"
	"google.golang.org/api/logging/v2"
	"google.golang.org/api/monitoring/v3"
	"google.golang.org/api/pubsub/v1"
	"google.golang.org/api/serviceusage/v1"
	"google.golang.org/api/storage/v1"

	computeBeta "google.golang.org/api/compute/v0.beta"
)

// CloudResourceManagerService returns the service connection for GCP Cloud Resource Manager service
func CloudResourceManagerService(ctx context.Context, connectionManager *connection.Manager) (*cloudresourcemanager.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "CloudResourceManagerService"
	if cachedData, ok := connectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*cloudresourcemanager.Service), nil
	}

	// so it was not in cache - create service
	svc, err := cloudresourcemanager.NewService(ctx)
	if err != nil {
		return nil, err
	}

	connectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// ComputeBetaService returns the service connection for GCP Compute service beta version
func ComputeBetaService(ctx context.Context, connectionManager *connection.Manager) (*computeBeta.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "ComputeBetaService"
	if cachedData, ok := connectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*computeBeta.Service), nil
	}

	// so it was not in cache - create service
	svc, err := computeBeta.NewService(ctx)
	if err != nil {
		return nil, err
	}

	connectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// ComputeService returns the service connection for GCP Compute service
func ComputeService(ctx context.Context, connectionManager *connection.Manager) (*compute.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "ComputeService"
	if cachedData, ok := connectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*compute.Service), nil
	}

	// so it was not in cache - create service
	svc, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}

	connectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// CloudFunctionsService returns the service connection for GCP Cloud Functions service
func CloudFunctionsService(ctx context.Context, connectionManager *connection.Manager) (*cloudfunctions.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "CloudFunctionsService"
	if cachedData, ok := connectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*cloudfunctions.Service), nil
	}

	// so it was not in cache - create service
	svc, err := cloudfunctions.NewService(ctx)
	if err != nil {
		return nil, err
	}

	connectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// IAMService returns the service connection for GCP IAM service
func IAMService(ctx context.Context, connectionManager *connection.Manager) (*iam.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "IAMService"
	if cachedData, ok := connectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*iam.Service), nil
	}

	// so it was not in cache - create service
	svc, err := iam.NewService(ctx)
	if err != nil {
		return nil, err
	}

	connectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// LoggingService returns the service connection for GCP Logging service
func LoggingService(ctx context.Context, connectionManager *connection.Manager) (*logging.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "LoggingService"
	if cachedData, ok := connectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*logging.Service), nil
	}

	// so it was not in cache - create service
	svc, err := logging.NewService(ctx)
	if err != nil {
		return nil, err
	}

	connectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// MonitoringService returns the service connection for GCP Monitoring service
func MonitoringService(ctx context.Context, connectionManager *connection.Manager) (*monitoring.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "MonitoringService"
	if cachedData, ok := connectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*monitoring.Service), nil
	}

	// so it was not in cache - create service
	svc, err := monitoring.NewService(ctx)
	if err != nil {
		return nil, err
	}

	connectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// PubsubService returns the service connection for GCP Pub/Sub service
func PubsubService(ctx context.Context, connectionManager *connection.Manager) (*pubsub.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "PubsubService"
	if cachedData, ok := connectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*pubsub.Service), nil
	}

	// so it was not in cache - create service
	svc, err := pubsub.NewService(ctx)
	if err != nil {
		return nil, err
	}

	connectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// ServiceUsageService returns the service connection for GCP Service Usage service
func ServiceUsageService(ctx context.Context, connectionManager *connection.Manager) (*serviceusage.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "ServiceUsageService"
	if cachedData, ok := connectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*serviceusage.Service), nil
	}

	// so it was not in cache - create service
	svc, err := serviceusage.NewService(ctx)
	if err != nil {
		return nil, err
	}

	connectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

// StorageService returns the service connection for GCP Storgae service
func StorageService(ctx context.Context, connectionManager *connection.Manager) (*storage.Service, error) {
	// have we already created and cached the service?
	serviceCacheKey := "StorageService"
	if cachedData, ok := connectionManager.Cache.Get(serviceCacheKey); ok {
		return cachedData.(*storage.Service), nil
	}

	// so it was not in cache - create service
	svc, err := storage.NewService(ctx)
	if err != nil {
		return nil, err
	}

	connectionManager.Cache.Set(serviceCacheKey, svc)
	return svc, nil
}

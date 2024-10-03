STEAMPIPE_INSTALL_DIR ?= ~/.steampipe
STEAMPIPE_PLUGIN_VERSION ?= latest
BUILD_TAGS = netgo
install:
	go build -o $(STEAMPIPE_INSTALL_DIR)/plugins/hub.steampipe.io/plugins/turbot/gcp@$(STEAMPIPE_PLUGIN_VERSION)/steampipe-plugin-gcp.plugin -tags "${BUILD_TAGS}" *.go

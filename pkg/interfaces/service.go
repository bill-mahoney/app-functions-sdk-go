//
// Copyright (c) 2021 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package interfaces

import (
	"net/http"

	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/command"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/coredata"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/logger"
	"github.com/edgexfoundry/go-mod-core-contracts/v2/clients/notifications"
	"github.com/edgexfoundry/go-mod-registry/v2/registry"

	bootstrapInterfaces "github.com/edgexfoundry/go-mod-bootstrap/v2/bootstrap/interfaces"
)

const (
	// AppServiceContextKey is the context key for getting the reference to the ApplicationService from the context passed to
	// a custom REST Handler
	AppServiceContextKey = "AppService"

	// ProfileSuffixPlaceholder is the placeholder text to use in an application service's service key if the
	// the name of the configuration profile used is to be used in the service's service key.
	// Only useful if the service has multiple configuration profiles to choose from at runtime.
	// Example:
	//    const (
	//		serviceKey = "MyServiceName-" + interfaces.ProfileSuffixPlaceholder
	//	  )
	ProfileSuffixPlaceholder = "<profile>"
)

// UpdatableConfig interface allows services to have custom configuration populated from configuration stored
// in the Configuration Provider (aka Consul). Services using custom configuration must implement this interface
// on their custom configuration, even if they do not use Configuration Provider. If they do not use the
// Configuration Provider they can have dummy implementation of this interface.
// This wraps the actual interface from go-mod-bootstrap so app service code doesn't have to have the additional
// direct import of go-mod-bootstrap.
type UpdatableConfig interface {
	bootstrapInterfaces.UpdatableConfig
}

// ApplicationService defines the interface for an edgex Application Service
type ApplicationService interface {
	// AddRoute a custom REST route to the application service's internal webserver
	// A reference to this ApplicationService is add the the context that is passed to the handler, which
	// can be retrieved using the `AppService` key
	AddRoute(route string, handler func(http.ResponseWriter, *http.Request), methods ...string) error
	// ApplicationSettings returns the key/value map of custom settings
	ApplicationSettings() map[string]string
	// GetAppSetting is a convenience function return a setting from the ApplicationSetting
	// section of the service configuration.
	// An error is returned if the specified setting is not found.
	GetAppSetting(setting string) (string, error)
	// GetAppSettingStrings is a convenience function that parses the value for the specified custom
	// application setting as a comma separated list. It returns the list of strings.
	// An error is returned if the specified setting is not found.
	GetAppSettingStrings(setting string) ([]string, error)
	// SetFunctionsPipeline set the functions pipeline with the specified list of Application Functions.
	// Note that the functions are executed in the order provided in the list.
	// An error is returned if the list is empty.
	SetFunctionsPipeline(transforms ...AppFunction) error
	// MakeItRun starts the configured trigger to allow the functions pipeline to execute when the trigger
	// receives data and starts the internal webserver. This is a long running function which does not return until
	// the service is stopped or MakeItStop() is called.
	// An error is returned if the trigger can not be create or initialized or if the internal webserver
	// encounters an error.
	MakeItRun() error
	// MakeItStop stops the configured trigger so that the functions pipeline no longer executes.
	// An error is returned
	MakeItStop()
	// RegisterCustomTriggerFactory registers a trigger factory for a custom trigger to be used.
	RegisterCustomTriggerFactory(name string, factory func(TriggerConfig) (Trigger, error)) error
	// Adds and returns a BackgroundPublisher which is used to publish asynchronously to the Edgex MessageBus.
	// Not valid for use with the HTTP or External MQTT triggers
	AddBackgroundPublisher(capacity int) BackgroundPublisher
	// GetSecret returns the secret data from the secret store (secure or insecure) for the specified path.
	// An error is returned if the path is not found or any of the keys (if specified) are not found.
	// Omit keys if all secret data for the specified path is required.
	GetSecret(path string, keys ...string) (map[string]string, error)
	// StoreSecret stores the specified secret data into the secret store (secure only) for the specified path
	// An error is returned if:
	//   - Specified secret data is empty
	//   - Not using the secure secret store, i.e. not valid with InsecureSecrets configuration
	//   - Secure secret provider is not properly initialized
	//   - Connection issues with Secret Store service.
	StoreSecret(path string, secretData map[string]string) error // LoggingClient returns the Logger client
	LoggingClient() logger.LoggingClient
	// EventClient returns the Event client. Note if Core Data is not specified in the Clients configuration,
	// this will return nil.
	EventClient() coredata.EventClient
	// CommandClient returns the Command client. Note if Support Command is not specified in the Clients configuration,
	// this will return nil.
	CommandClient() command.CommandClient
	// NotificationsClient returns the Notifications client. Note if Support Notifications is not specified in the
	// Clients configuration, this will return nil.
	NotificationsClient() notifications.NotificationsClient
	// RegistryClient() returns the Registry client. Note the registry must been enable, otherwise this will return nil.
	// Useful if service needs to add additional health checks or needs to get endpoint of another registered service
	RegistryClient() registry.Client
	// LoadConfigurablePipeline loads the function pipeline from configuration.
	// An error is returned if the configuration is not valid, i.e. missing required function parameters,
	// invalid function name, etc.
	// Only useful if pipeline from configuration is always defined in configuration as in App Service Configurable.
	LoadConfigurablePipeline() ([]AppFunction, error)
	// LoadCustomConfig loads the service's custom configuration from local file or the Configuration Provider (if enabled)
	// Configuration Provider will also be seeded with the custom configuration if service is using the Configuration Provider.
	// UpdateFromRaw interface will be called on the custom configuration when the configuration is loaded from the
	// Configuration Provider.
	LoadCustomConfig(config UpdatableConfig, sectionName string) error
	// ListenForCustomConfigChanges starts a listener on the Configuration Provider for changes to the specified
	// section of the custom configuration. When changes are received from the Configuration Provider the
	// UpdateWritableFromRaw interface will be called on the custom configuration to apply the updates and then signal
	// that the changes occurred via writableChanged.
	ListenForCustomConfigChanges(configToWatch interface{}, sectionName string, changedCallback func(interface{})) error
}

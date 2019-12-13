package openc

type (
	// Service is a Service a pipeline should execute
	Service struct {
		// Name is official service which is part of the service catalog (https://github.com/open-integration/core-services/releases)
		Name string
		// Version of the service, empty string will use the latest version from catalog
		Version string
		// Path a location of the local fs the service can be found, Path cannot be set with Name together
		Path string
		// As alias name to refer the service as part of the task implementation
		As string
	}
)

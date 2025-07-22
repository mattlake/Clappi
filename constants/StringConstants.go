package constants

const (
	// Panel Titles
	EnvironmentsPanelTitle = "Environments"
	ApisPanelTitle         = "APIs"
	EndpointsPanelTitle    = "Endpoints"
	RequestPanelTitle      = "Request"
	ResponsePanelTitle     = "Response"

	// Errors
	ApiLoadingError         = "could not load APIs: %s"
	SpecLoadingError        = "could not load spec: %s"
	FileAccessError         = "could not access path %s: %w"
	FileReadingError        = "could not read file %s: %w"
	FileParsingError        = "could not parse file %s: %w"
	UnsupportedVersionError = "v2 is not currently supported. Please use a V3 specification"
	ModelBuildingError      = "could not build model from file %s"
)

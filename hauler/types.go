package hauler

type RegistryTarget struct {
	Url      string
	Username string
	Password string
	Image    string
	Tag      string
}

type Hauler struct {
	StorageDir                string
	SourceRegistryTarget      RegistryTarget
	DestinationRegistryTarget RegistryTarget
}

type Config struct {
	SourceRegistry      string
	SourceImage         string
	SourceTag           string
	SourceUsername      string
	SourcePassword      string
	DestinationRegistry string
	DestinationImage    string
	DestinationTag      string
	DestinationUsername string
	DestinationPassword string
	StorageDir          string
}

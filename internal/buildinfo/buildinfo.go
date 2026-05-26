package buildinfo

var Version = "dev"
var BuildDate = "unknown"

func GetBuildInfo() string {
	return "Version: " + Version + ", Build Date: " + BuildDate
}

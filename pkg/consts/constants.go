package consts

const Version = "0.0.1-alpha"

// BasePath defines the default path, overridden during production builds.
var BasePath string = "./false_location"

// Testing mode flag
var Testing bool = false

// GetBasePath retrieves the effective base path
func GetBasePath() string {
	return BasePath
}

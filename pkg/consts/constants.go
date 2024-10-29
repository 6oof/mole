package consts

const Version = "0.0.1-alpha"

// This gets changed if env var MOLE_ENV_PROD is set to 1 (look at main.go).
// It should never be changed anywhere else.
var BasePath string = "./false_location"

var Testing bool = false

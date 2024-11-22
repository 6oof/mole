package cmd

var (
	domainFlag      string
	portFlag        int
	locationFlag    string
	repositoryFlag  string
	descriptionFlag string
	branchFlag      string
	confirmFlag     bool
	hardRerloadFlag bool
	deployDown      bool
)

// flags for service actions
var serviceStartFlag, serviceStopFlag, serviceEnableFlag, serviceDisableFlag bool

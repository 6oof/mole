package cmd

var (
	domainFlag      string
	portFlag        string
	locationFlag    string
	repositoryFlag  string
	descriptionFlag string
	branchFlag      string
	confirmFlag     bool
	pTypeFlag       string
	hardRerloadFlag bool
)

// flags for service actions
var serviceStartFlag, serviceStopFlag, serviceEnableFlag, serviceDisableFlag bool

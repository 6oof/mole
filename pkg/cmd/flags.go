package cmd

var (
	domainFlag      string
	portFlag        int
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

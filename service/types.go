package service

// SwarmService defines the service and their tasks
type SwarmService struct {
	ID   string
	Name string
	Task []SwarmTask
}

// SwarmTask defines the tasks running in a cluster
type SwarmTask struct {
	ID      string
	Name    string
	Address string
	Port    int
	Checks  []Check
}

// Check defines the health check for the application
type Check struct {
	ID                       string
	Name                     string
	Interval                 string
	Timeout                  string
	HTTP                     string
	Path                     string
	Header                   map[string][]string
	RemoveFailedServiceAfter string
}

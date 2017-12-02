package models

type (
	// FlavorsData struct represents the values of List Flavor API
	FlavorsData struct {
		Name  string
		ID    string
		RAM   int
		VCPUs int
		Disk  int
		RXTX  float64
	}
	// ImagesData struct represents the values of List Images API
	ImagesData struct {
		Name     string
		ID       string
		MinDisk  int
		Status   string
		Progress int
		MinRAM   int
		Metadata map[string]string
	}
	// InstancesData struct represents the values of List Instances API
	InstancesData struct {
		Name          string
		ID            string
		Owner         string
		Image         map[string]interface{}
		Flavor        map[string]interface{}
		Host          string
		Status        string
		SecurityGroup []map[string]interface{}
		CreatedAt     string
		UpdatedAt     string
	}
)

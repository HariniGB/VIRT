package models

import "strings"

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
		Type     string
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

	Stack struct {
		Name string `json:"stack_name"`
		Status string `json:"stack_status"`
		CreationTime string `json:"creation_time"`
		Outputs []StackOutput `json:"outputs"`
	}

	QuotaEntry struct {
		Name string `json:"Name"`
		Value int `json:"Value"`
	}

	StackOutput struct {
		Key string `json:"output_key"`
		Value string `json:"output_value"`
		Description string `json:"description"`
	}

	StackResponse struct {
		Name string `json:"name"`
		Status string `json:"status"`
		CreationTime string `json:"creation_time"`
		PrivateIp string `json:"private_ip"`
		PublicIp string `json:"public_ip"`
	}

	ProvisionRequest struct {
		Name string `json:"name"`
		Flavor string `json:"flavor"`
		Type string `json:"type"`
	}

	QuotaResponse struct {
		CpuUsage      int `json:"cpu_usage"`
		CpuMax        int `json:"cpu_max"`
		MemoryUsage   int `json:"memory_usage"`
		MemoryMax     int `json:"memory_max"`
		InstanceUsage int `json:"instance_usage"`
		InstanceMax   int `json:"instance_max"`
	}
)

func (s *Stack) ToStackResponse() *StackResponse {
	if s.Name == "" || strings.Index(s.Name, "-network") != -1 {
		return nil
	}

	resp := &StackResponse{
		Name: s.Name,
		Status: s.Status,
		CreationTime: s.CreationTime,
	}

	if s.Outputs != nil {
		for _, out := range s.Outputs {
			if out.Key == "server1_private_ip" {
				resp.PrivateIp = out.Value
			} else if out.Key == "server1_public_ip" {
				resp.PublicIp = out.Value
			}
		}
	}
	return resp
}

func (q *QuotaResponse) FromQuotas(quotas []QuotaEntry) {
	for _, quota := range quotas {
		switch quota.Name {
		case "maxTotalInstances":
			q.InstanceMax = quota.Value
		case "totalInstancesUsed":
			q.InstanceUsage = quota.Value
		case "maxTotalCores":
			q.CpuMax = quota.Value
		case "totalCoresUsed":
			q.CpuUsage = quota.Value
		case "maxTotalRAMSize":
			q.MemoryMax = quota.Value
		case "totalRAMUsed":
			q.MemoryUsage = quota.Value
		}
	}
}

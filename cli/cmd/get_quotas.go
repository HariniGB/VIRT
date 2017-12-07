package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"github.com/HariniGB/openstack-api/controllers"
	"fmt"
)

var (
	quota_table_headers = []string {"Name", "Current", "Total"}
)

func NewGetQuotasCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "quotas",
		Short: "get all quotas",
		Run: func(cmd *cobra.Command, args []string) {
			if err := RunGetQuotas(cmd, args); err != nil {
				log.Println(err)
			}
		},
	}
	return cmd
}

func RunGetQuotas(cmd *cobra.Command, args []string) error{

	quotas := controllers.GetQuotas()
	if quotas == nil {
		return fmt.Errorf("unable to get quotas")
	}

	data := [][]string {
		{"CPU", fmt.Sprint(quotas.CpuUsage), fmt.Sprint(quotas.CpuMax)},
		{"Memory", fmt.Sprint(quotas.MemoryUsage), fmt.Sprint(quotas.MemoryMax)},
		{"Instances", fmt.Sprint(quotas.InstanceUsage), fmt.Sprint(quotas.InstanceMax)},
	}

	newTable(quota_table_headers, data)
	return nil
}

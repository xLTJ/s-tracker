package cmd

import "github.com/spf13/cobra"

const (
	applicationApiPath = "https://mit.s.dk/api/building"
)

var applicationsCmd = &cobra.Command{
	Use:   "applications",
	Short: "Get a list of all your applications",
	RunE: func(cmd *cobra.Command, args []string) error {
		//client.Collector.Visit()
		return nil
	},
}

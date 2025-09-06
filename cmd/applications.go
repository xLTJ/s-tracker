package cmd

import (
	"fmt"
	"github.com/cheynewallace/tabby"
	"github.com/spf13/cobra"
)

var applicationsCmd = &cobra.Command{
	Use:   "applications",
	Short: "Get a list of all your applications",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("\nGetting buildings you have applications for...\n\n")
		buildings, err := client.GetUserBuildings()
		if err != nil {
			return fmt.Errorf("error getting buildings: %v", err)
		}

		table := tabby.New()
		table.AddHeader("Id", "Name", "Address", "Municipality")
		for _, building := range buildings {
			table.AddLine(
				building.Id,
				building.Name,
				building.Address,
				building.Municipality,
			)
		}

		table.Print()
		return nil
	},
}

package cmd

import (
	"fmt"
	"github.com/cheynewallace/tabby"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	fullInfo        bool
	applicationsCmd = &cobra.Command{
		Use:   "applications",
		Short: "Get a list of all your applications",
		RunE: func(cmd *cobra.Command, args []string) error {
			pterm.Info.Printf("Getting buildings you have applications for...")
			pterm.Println()
			buildings, err := client.GetUserBuildings()
			if err != nil {
				return fmt.Errorf("error getting buildings: %v", err)
			}

			table := tabby.New()

			if fullInfo {
				buildingsFull, err := client.GetUserBuildingsFull(buildings)
				if err != nil {
					return fmt.Errorf("error getting full building info: %v", err)
				}

				table.AddHeader("Id", "Name", "Address", "Municipality", "Build Year", "Best Waiting List Category")
				for _, building := range buildingsFull {
					table.AddLine(
						building.BasicBuildingInfo.Id,
						building.BasicBuildingInfo.Name,
						building.BasicBuildingInfo.Address,
						building.BasicBuildingInfo.Municipality,
						building.BuildYear,
						colorCategoryHelper(building.BestWaitingListCategory),
					)
				}
			} else {
				table.AddHeader("Id", "Name", "Address", "Municipality")
				for _, building := range buildings {
					table.AddLine(
						building.Id,
						building.Name,
						building.Address,
						building.Municipality,
					)
				}
			}

			table.Print()
			return nil
		},
	}
)

func colorCategoryHelper(category rune) (coloredString string) {
	switch category {
	case 'A':
		coloredString = pterm.FgBlue.Sprint(string(category))
	case 'B':
		coloredString = pterm.FgGreen.Sprint(string(category))
	case 'C':
		coloredString = pterm.FgYellow.Sprint(string(category))
	default:
		coloredString = pterm.FgRed.Sprint(string(category))

	}
	return
}

func init() {
	applicationsCmd.Flags().BoolVarP(&fullInfo, "all", "a", false, "List all info about applications. Will take longer")
}

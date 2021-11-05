package main

import (
	"os"

	"github.com/kevholditch/vacuum/internal/app/vacuum"
	"github.com/liamg/tml"
	"github.com/spf13/cobra"
)

var regions []string

var defaultRegions = []string{"eu-west-1", "eu-west-2"}

var volumesCmd = &cobra.Command{
	Use:   "volumes",
	Short: "remove available EC2 Volumes that are not being used",
	Run: func(cmd *cobra.Command, args []string) {

		for _, region := range regions {
			r := vacuum.Region(region)
			tml.Printf("[%s]\n", region)
			tml.Printf("\t Identifying volumes\n")

			resources, err := vacuum.Volumes().Identify(r)
			if err != nil {
				tml.Printf("<bold><red>Error:</red></bold> could not check volumes error: %s\n", err)
				os.Exit(1)
			}
			tml.Printf("\t Found %d volume(s)\n", len(resources.Resources()))

			if len(resources.Resources()) > 0 {
				tml.Printf("\t Vacuuming %d volumes(s)...\n", len(resources.Resources()))
				err = vacuum.Volumes().Clean(resources)
				tml.Printf("\t Done.\n")
			}

			tml.Printf("\n")

		}

	},
}

func init() {
	volumesCmd.Flags().StringSliceVarP(&regions, "regions", "r", regions, "AWS regions you want to clean e.g. eu-west-1")

	rootCmd.AddCommand(volumesCmd)
}

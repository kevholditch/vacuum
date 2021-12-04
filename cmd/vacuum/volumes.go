package main

import (
	"github.com/kevholditch/vacuum/internal/app/vacuum"
	"github.com/spf13/cobra"
)

var volumesCmd = &cobra.Command{
	Use:   "volumes",
	Short: "remove available EC2 Volumes that are not being used",
	Run: func(cmd *cobra.Command, args []string) {
		vacuum.Vacuum(regions, vacuum.Volumes())
	},
}

func init() {
	volumesCmd.Flags().StringSliceVarP(&regions, "regions", "r", regions, "AWS regions you want to clean e.g. eu-west-1")
	rootCmd.AddCommand(volumesCmd)
}

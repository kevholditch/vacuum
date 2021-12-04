package main

import (
	"github.com/kevholditch/vacuum/internal/app/vacuum"
	"github.com/spf13/cobra"
)

var eniCmd = &cobra.Command{
	Use:   "enis",
	Short: "remove available ENIs that are not being used",
	Run: func(cmd *cobra.Command, args []string) {
		vacuum.Vacuum(regions, vacuum.NetworkInterfaces())
	},
}

func init() {
	eniCmd.Flags().StringSliceVarP(&regions, "regions", "r", regions, "AWS regions you want to clean e.g. eu-west-1")
	rootCmd.AddCommand(eniCmd)
}

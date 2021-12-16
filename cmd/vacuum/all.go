package main

import (
	"github.com/kevholditch/vacuum/internal/app/vacuum"
	"github.com/spf13/cobra"
)

var allCmd = &cobra.Command{
	Use:   "all",
	Short: "remove dirt from your AWS account",
	Long:  "remove available unused resources from your AWS account, stop lining Jeff Bezos's pockets!",
	Run: func(cmd *cobra.Command, args []string) {
		vacuum.Vacuum(regions, vacuum.Volumes(), vacuum.NetworkInterfaces(), vacuum.SecurityRules())
	},
}

func init() {
	allCmd.Flags().StringSliceVarP(&regions, "regions", "r", regions, "AWS regions you want to clean e.g. eu-west-1")
	rootCmd.AddCommand(allCmd)
}

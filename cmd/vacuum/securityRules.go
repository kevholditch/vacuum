package main

import (
	"github.com/kevholditch/vacuum/internal/app/vacuum"
	"github.com/spf13/cobra"
)

// securityRulesCmd represents the securityRules command
var securityRulesCmd = &cobra.Command{
	Use:   "securityRules",
	Short: "remove security rules created by K8s annotations (via aws-load-balancer-controller) which reference a removed NLB",
	Run: func(cmd *cobra.Command, args []string) {
		vacuum.Vacuum(regions, vacuum.SecurityRules())
	},
}

func init() {
	securityRulesCmd.Flags().StringSliceVarP(&regions, "regions", "r", regions, "AWS regions you want to clean e.g. eu-west-1")
	rootCmd.AddCommand(securityRulesCmd)
}

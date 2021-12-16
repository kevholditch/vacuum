package main

import (
	"fmt"

	"github.com/kevholditch/vacuum/internal/app/vacuum"

	"github.com/spf13/cobra"
)

var regions []string
var defaultRegions = []string{"eu-west-1", "eu-west-2"}

var rootCmd = &cobra.Command{
	Use:   "vacuum",
	Short: "Clean your AWS account of unused resources",
	Long:  `Vacuum will deep clean your AWS account, saving you money!!`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		fmt.Printf(`
                             ▒▒▒▒▒▒▒▒▒▒▒▒                    
                          ▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒                  
        ░░░░░░░░░░        ▒▒▒▒        ▒▒▒▒▒▒                
      ░░░░░░░░░░░░░░      ▒▒▒▒          ▒▒▒▒                
      ▒▒▒▒▒▒▒▒▒▒▒▒▒▒      ▒▒▒▒▒▒        ▒▒▒▒                
  ░░░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░░░    ▓▓▒▒▒▒      ▒▒▒▒                
░░░░░░░░▒▒▒▒▒▒▒▒▒▒░░░░░░░░    ▒▒▒▒      ▒▒▒▒                
▒▒░░░░░░░░░░░░░░░░░░░░░░▒▒    ▒▒▒▒        ▒▒▒▒              
▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒    ▒▒▒▒        ▒▒▒▒              
░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░    ▒▒▒▒        ▒▒▒▒              
░░░░░░▒▒▒▒▒▒▒▒▒▒▒▒▒▒░░░░░░  ▒▒▒▒▒▒        ▒▒▒▒              
░░░░░░░░░░░░░░░░░░░░░░░░▒▒▒▒▒▒▒▒          ▒▒▒▒              
░░░░░░░░░░░░░░░░░░░░░░░░▒▒▒▒▒▒            ▒▒▒▒              
░░░░░░░░░░░░░░░░░░░░░░░░░░                ▒▒▒▒              
░░░░░░░░░░░░░░░░░░░░░░░░░░                  ▒▒▓▓            
░░░░░░░░░░░░░░░░░░░░░░░░░░                  ▒▒▒▒            
░░░░░░░░░░░░░░░░░░░░░░░░░░                  ▒▒▒▒            
▒▒░░░░░░░░░░░░░░░░░░░░░░▒▒                  ▒▒▒▒            
▒▒▒▒▒▒▒▒░░░░░░░░░░▒▒▒▒▒▒▒▒                  ▒▒▒▒            
██▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒▒██                ▒▒▒▒▒▒▒▒          
████    ▒▒▒▒▒▒▒▒▒▒    ████              ▒▒▒▒▒▒▒▒▒▒▒▒        
          ████                            ░░░░░░░░░░░░░░░░░░
                                        ░░░░░░░░░░░░░░░░░░  
                                            ░░░░░░░░░░    

https://github.com/kevholditch/vacuum %s

`, vacuum.Version())
	},
}

func init() {
	for _, region := range defaultRegions {
		regions = append(regions, region)
	}
}

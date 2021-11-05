package main

import (
	"fmt"

	"github.com/kevholditch/vacuum/internal/app/vacuum"

	"github.com/spf13/cobra"
)

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

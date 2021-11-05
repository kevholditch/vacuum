package vacuum

import (
	"fmt"
	"os"

	"github.com/liamg/tml"
)

func clearLine() {
	fmt.Printf("\033[2K\r")
}

func wrapWhite(str string) string {
	return fmt.Sprintf("<white>%s</white>", str)
}

func Vacuum(regions []string, vacuumers ...Vacuumer) {

	for _, region := range regions {
		r := Region(region)
		for _, v := range vacuumers {
			tml.Printf("<white>[%s]</white>\n", region)
			tml.Printf("\t [%s] Identifying...", v.Type())

			resources, err := v.Identify(r)
			if err != nil {
				tml.Printf("<bold><red>Error:</red></bold> could not identify %s, details: %s\n", v.Type(), err)
				os.Exit(1)
			}
			clearLine()
			foundStr := fmt.Sprintf("\t [%s] Found %d\n", v.Type(), len(resources.Resources()))
			if len(resources.Resources()) > 0 {
				foundStr = wrapWhite(foundStr)
			}
			tml.Printf(foundStr)

			if len(resources.Resources()) > 0 {
				tml.Printf("\t Vacuuming %d %s...\n", len(resources.Resources()), v.Type())
				err = v.Clean(resources)
				tml.Printf("\t Done.\n")
			}

			tml.Printf("\n")
		}

	}
}

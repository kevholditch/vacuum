package vacuum

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/liamg/tml"
)

func clearLine() {
	fmt.Printf("\033[2K\r")
}

func wrapWhite(str string) string {
	return fmt.Sprintf("<white>%s</white>", str)
}

func Vacuum(regions []string, vacuumers ...Vacuumer) {

	reader := bufio.NewReader(os.Stdin)

	for _, region := range regions {
		r := Region(region)
		tml.Printf("\n<white>[%s]</white>\n", region)

		var clean []struct {
			vaccuumer Vacuumer
			resources Resources
		}

		for _, v := range vacuumers {
			tml.Printf("\t [%s] Identifying...", v.Type())

			resources, err := v.Identify(r)
			if err != nil {
				tml.Printf("<bold><red>Error:</red></bold> could not identify %s, details: %s\n", v.Type(), err)
				os.Exit(1)
			}
			clearLine()
			foundStr := fmt.Sprintf("\t [%s] Found %d", v.Type(), len(resources.Resources()))
			if len(resources.Resources()) > 0 {
				foundStr = wrapWhite(foundStr)
			}
			tml.Printf(foundStr)

			if len(resources.Resources()) > 0 {
				clean = append(clean, struct {
					vaccuumer Vacuumer
					resources Resources
				}{vaccuumer: v, resources: resources})

			}

			tml.Printf("\n")
		}

		if len(clean) > 0 {
			tml.Printf("\nVacuum? [y/n]:")
			now := time.Now()
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)
			if strings.EqualFold("y", text) {

				for _, c := range clean {
					c.vaccuumer.Clean(c.resources, func(amount int) {
						clearLine()
						tml.Printf("Vacuuming %s [%d/%d]", c.vaccuumer.Type(), amount+1, len(c.resources.Resources()))
					})
				}
				clearLine()

				tml.Printf("Vacuumed in %.2fs", time.Now().Sub(now).Seconds())
				tml.Printf("\n\n")
			} else {
				tml.Printf("no worries, will leave it dirty!\n\n")
			}
		}

	}
}

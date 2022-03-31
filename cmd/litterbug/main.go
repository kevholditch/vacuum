package main

import "github.com/kevholditch/vacuum/internal/app/litterbug"

/*
Create some litterbug in your AWS account for vacuum to clean
*/
func main() {
	bug, err := litterbug.InRegion("eu-west-1")
	if err != nil {
		panic(err)
	}

	err = bug.CreateAvailableVolumes(100).Litter()
	if err != nil {
		panic(err)
	}
}

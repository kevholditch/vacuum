package app

import (
	"testing"
)

func TestRemoveAvailableVolumeInRegion(t *testing.T) {

	given, when, then, _ := newVolumesTest(t)

	given.
		an_available_volume_exists_in_region("eu-west-1")

	when.
		volumes_are_vacuumed_in("eu-west-1")

	then.
		there_should_be_no_available_volumes_in_the_region("eu-west-1")

}

func TestLeaveVolumeInOtherRegion(t *testing.T) {

	given, when, then, clean := newVolumesTest(t)
	defer clean()

	given.
		an_available_volume_exists_in_region("eu-west-2")

	when.
		volumes_are_vacuumed_in("eu-west-1")

	then.
		there_should_be_no_available_volumes_in_the_region("eu-west-1")

}

func TestIdentifyVolumes(t *testing.T) {
	given, when, then, clean := newVolumesTest(t)
	defer clean()

	given.
		three_available_volumes_exist_in("eu-west-1")

	when.
		the_available_volumes_are_identified("eu-west-1")

	then.
		there_should_be_three_available_volumes_identified_in("eu-west-1")

}

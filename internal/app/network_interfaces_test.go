package app

import "testing"

func TestRemoveAvailableNetworkInterfaceInRegion(t *testing.T) {

	given, when, then, clean := newNetworkInterfacesTest(t)
	defer clean()

	given.
		an_available_network_interface_exists_in_region("eu-west-1")

	when.
		network_interfaces_are_vacuumed_in("eu-west-1")

	then.
		there_should_be_no_available_network_interfaces_in_the_region("eu-west-1")

}

func TestLeaveNetworkInterfaceInDifferentRegion(t *testing.T) {

	given, when, then, clean := newNetworkInterfacesTest(t)
	defer clean()

	given.
		an_available_network_interface_exists_in_region("eu-west-2")

	when.
		network_interfaces_are_vacuumed_in("eu-west-1")

	then.
		there_should_be_no_available_network_interfaces_in_the_region("eu-west-1")
}

func TestIdentifyNetworkInterfaces(t *testing.T) {
	given, when, then, clean := newNetworkInterfacesTest(t)
	defer clean()

	given.
		three_available_network_interfaces_exist_in("eu-west-1")

	when.
		the_available_network_interfaces_are_identified("eu-west-1")

	then.
		there_should_be_three_available_network_interfaces_identified_in("eu-west-1")

}
